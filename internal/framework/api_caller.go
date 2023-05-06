package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const ApiCallerWrongRequest ErrorIdentifier = "ApiCallerWrongRequest"
const ApiCallerAuthenticationRequired ErrorIdentifier = "ApiCallerAuthenticationIsRequired"
const ApiCallerWrongResponse ErrorIdentifier = "ApiCallerWrongResponse"
const ApiCallerNotSupportedMethod ErrorIdentifier = "ApiCallerNotSupportedMethod"
const ApiCallerNetworkIssues ErrorIdentifier = "ApiCallerNetworkIssues"

type File struct {
	FileName  string
	Reader    io.Reader
	FieldName string
}

type RetryConfig struct {
	MaxRetires      int
	MaxWaitInSecond int
	SleepTimeUnit   time.Duration
}

func NewRetryConfig(maxRetires int, maxWaitInSecond int, sleepTimeUnit time.Duration) *RetryConfig {
	return &RetryConfig{MaxRetires: maxRetires, MaxWaitInSecond: maxWaitInSecond, SleepTimeUnit: sleepTimeUnit}
}

type ApiCaller[T any] struct {
	logger     *ApiCallerLogger
	CheckRetry func(ctx context.Context, resp *http.Response, err error) (bool, error)
	httpClient *http.Client
}

func NewApiCaller[T any](logger *ApiCallerLogger) *ApiCaller[T] {
	return &ApiCaller[T]{
		logger:     logger,
		CheckRetry: retryablehttp.DefaultRetryPolicy,
	}
}

func (a *ApiCaller[T]) MakeRequest(
	ctx context.Context,
	method string,
	url string,
	params map[string]interface{},
	retryConfig RetryConfig,
	additionalHeaders *http.Header,
) (T, error) {
	var result T
	req, err := a.prepareJsonRequest(ctx, method, url, params)
	if err != nil {
		return result, err
	}

	if additionalHeaders != nil {
		for name, values := range *additionalHeaders {
			for _, value := range values {
				req.Header.Set(name, value)
			}
		}
	}

	resp, err := a.processRequest(
		ctx,
		req,
		retryConfig.MaxRetires,
		retryConfig.MaxWaitInSecond,
		retryConfig.SleepTimeUnit,
	)
	if err != nil {
		return result, err
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		httpError := NewTranslatedError(
			ctx,
			ApiCallerAuthenticationRequired,
			"Authentication is required for the remote API %s",
			url,
		)
		return result, httpError
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		httpError := NewTranslatedError(
			ctx,
			ApiCallerWrongRequest,
			"Request to the remote API %s is wrong: %s. Status code: %d",
			url,
			params,
			resp.StatusCode,
		)
		return result, httpError
	}
	if resp.StatusCode > 500 {
		httpError := NewTranslatedError(
			ctx,
			ApiCallerWrongResponse,
			"The remote API %s is unavailable for the getting with parameters %s.",
			url,
			params,
		)
		return result, httpError
	}
	errJson := json.Unmarshal(bodyBytes, &result)
	if errJson != nil {
		a.logger.Error("Response cannot be serialized", zap.Error(errJson), zap.String("response", string(bodyBytes)))
		return result, NewTranslatedError(
			ctx,
			ApiCallerWrongResponse,
			"The remote API %s returns invalid response.",
			url,
			params,
		)
	}
	return result, nil
}

func (a *ApiCaller[T]) UploadFile(
	ctx context.Context,
	url string,
	params map[string]string,
	file File,
	retryConfig RetryConfig,
) (T, error) {
	var result T
	req, err := a.prepareMultipartRequest(url, params, file)
	if err != nil {
		return result, err
	}

	resp, err := a.processRequest(
		ctx,
		req,
		retryConfig.MaxRetires,
		retryConfig.MaxWaitInSecond,
		retryConfig.SleepTimeUnit,
	)
	if err != nil {
		return result, err
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		httpError := NewTranslatedError(
			ctx,
			ApiCallerWrongRequest,
			"Request to the remote API %s is wrong: %s. Status code: %d",
			url,
			params,
			resp.StatusCode,
		)
		return result, httpError
	}
	errJson := json.Unmarshal(bodyBytes, &result)
	if errJson != nil {
		a.logger.Error("Response cannot be serialized", zap.Error(errJson), zap.String("response", string(bodyBytes)))
		return result, NewTranslatedError(
			ctx,
			ApiCallerWrongResponse,
			"The remote API %s returns invalid response.",
			url,
			params,
		)
	}
	return result, nil
}

func (a ApiCaller[T]) prepareJsonRequest(
	ctx context.Context,
	method string,
	url string,
	params map[string]interface{},
) (*http.Request, error) {
	var jsonReq []byte
	var err error
	if params == nil {
		jsonReq = []byte{}
	} else {
		jsonReq, err = json.Marshal(params)
		if err != nil {
			a.logger.Error("Wrong preparation request json.marshal", zap.Error(err))
			return nil, err
		}
	}

	var req *http.Request
	// Create request
	if method == "GET" {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			a.logger.Error("Wrong preparation request", zap.Error(err))
			return nil, err
		}
		q := req.URL.Query()
		for key, val := range params {
			if valStr, ok := val.(string); ok {
				q.Add(key, valStr)
			}
		}
		req.URL.RawQuery = q.Encode()
	} else {
		req, err = http.NewRequest(
			method,
			url, bytes.NewBuffer(jsonReq),
		)
		if err != nil {
			a.logger.Error("Wrong preparation request", zap.Error(err))
			return nil, err
		}
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

func (a ApiCaller[T]) prepareMultipartRequest(
	urlString string,
	params map[string]string,
	file File,
) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(file.FieldName, file.FileName)

	_, err := io.Copy(part, file.Reader)
	if err != nil {
		a.logger.Error("Wrong copying", zap.Error(err))
		return nil, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	writer.Close()

	urlObject, err := NewUrl(urlString)
	if err != nil {
		a.logger.Error("Cannot parse url", zap.Error(err))
		return nil, err
	}
	req, err := http.NewRequest("POST", urlObject.GetUrlWithoutBasicAuth(), body)
	if err != nil {
		a.logger.Error("Wrong preparation request", zap.Error(err))
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	if urlObject.HasBasicAuth() {
		username, password := urlObject.GetBasicAuth()
		req.SetBasicAuth(username, password)
	}

	return req, nil
}

func (a *ApiCaller[T]) processRequest(
	ctx context.Context,
	req *http.Request,
	maxRetires int,
	maxWaitInSecond int,
	sleepTimeUnit time.Duration,
) (*http.Response, error) {
	client := retryablehttp.NewClient()
	client.HTTPClient = a.getClient()
	client.RetryWaitMax = time.Second * time.Duration(maxWaitInSecond)
	client.RetryWaitMin = sleepTimeUnit
	client.RetryMax = maxRetires
	client.Logger = a.logger
	client.CheckRetry = a.CheckRetry
	req = req.WithContext(ctx)

	retryableRequest, err := retryablehttp.FromRequest(req)
	if err != nil {
		a.logger.Error("Wrong request", zap.Error(err))
		return nil, err
	}
	resp, err := client.Do(retryableRequest)
	if err != nil {
		a.logger.Error("Network issues", zap.Error(err))
		return nil, NewTranslatedError(ctx, ApiCallerNetworkIssues, "Network issues of %s", req.URL.String())
	}

	return resp, nil
}

func (a *ApiCaller[T]) getClient() *http.Client {
	if a.httpClient == nil {
		a.httpClient = cleanhttp.DefaultPooledClient()
	}
	a.httpClient.Timeout = time.Second * 10
	return a.httpClient
}

func (a *ApiCaller[T]) MakeBatchRequest(
	ctx context.Context,
	method string,
	url string,
	params []map[string]interface{},
	retryConfig RetryConfig,
) ([]byte, error) {
	req, err := a.prepareBatchJsonRequest(ctx, method, url, params)
	if err != nil {
		return []byte{}, err
	}

	resp, err := a.processRequest(
		ctx,
		req,
		retryConfig.MaxRetires,
		retryConfig.MaxWaitInSecond,
		retryConfig.SleepTimeUnit,
	)
	if err != nil {
		return []byte{}, err
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		httpError := NewTranslatedError(
			ctx,
			ApiCallerAuthenticationRequired,
			"Authentication is required for the remote API %s",
			url,
		)
		return []byte{}, httpError
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		httpError := NewTranslatedError(
			ctx,
			ApiCallerWrongRequest,
			"Request to the remote API %s is wrong: %s. Status code: %d",
			url,
			params,
			resp.StatusCode,
		)
		return []byte{}, httpError
	}
	if resp.StatusCode > 500 {
		httpError := NewTranslatedError(
			ctx,
			ApiCallerWrongResponse,
			"The remote API %s is unavailable for the getting with parameters %s.",
			url,
			params,
		)
		return []byte{}, httpError
	}
	return bodyBytes, nil
}

func (a *ApiCaller[T]) prepareBatchJsonRequest(
	ctx context.Context,
	method string,
	url string,
	params []map[string]interface{},
) (*http.Request, error) {
	jsonReq, err := json.Marshal(params)
	if err != nil {
		a.logger.Error("Wrong preparation request json.marshal", zap.Error(err))
		return nil, err
	}

	var req *http.Request
	// Create request
	if method == "POST" {
		req, err = http.NewRequest(
			method,
			url, bytes.NewBuffer(jsonReq),
		)
		if err != nil {
			a.logger.Error("Wrong preparation request", zap.Error(err))
			return nil, err
		}
		req.Header.Add("Content-Type", "application/json")
	} else {
		return nil, NewTranslatedError(ctx, ApiCallerNotSupportedMethod, "Method %s is not supported", method)
	}

	return req, nil
}
