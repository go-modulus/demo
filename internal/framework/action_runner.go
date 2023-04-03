package framework

import (
	"context"
	"encoding/json"
	"github.com/pasztorpisti/qs"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

var qsErrRegexp = regexp.MustCompile(`entry "([^"]+)" :: ([^:]+)`)

type ActionRunner struct {
	logger     Logger
	jsonWriter JsonResponseWriter
	router     Router
}

type ActionResponse struct {
	StatusCode      int
	Response        any
	Error           *ActionError
	IsLoggingErrors bool
	Template        *Page
}

func NewSuccessResponse(response any) ActionResponse {
	return ActionResponse{
		StatusCode: 200,
		Response:   response,
	}
}

func NewSuccessHtmlResponse(response any, template *Page) ActionResponse {
	return ActionResponse{
		StatusCode: 200,
		Response:   response,
		Template:   template,
	}
}

func NewSuccessCreationResponse(response any) ActionResponse {
	return ActionResponse{
		StatusCode: 201,
		Response:   response,
	}
}

func NewActionRunner(logger Logger, jsonWriter JsonResponseWriter, router Router) *ActionRunner {
	return &ActionRunner{logger: logger, jsonWriter: jsonWriter, router: router}
}

func (j *ActionRunner) Run(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, request any) ActionResponse,
	request any,
) {
	switch r.Method {
	case http.MethodGet, http.MethodDelete:
		j.runGet(w, r, action, request)
	case http.MethodPost:
		j.runPost(w, r, action, request)
	case http.MethodPut, http.MethodPatch:
		j.runPut(w, r, action, request)
	default:
		j.logger.Error(r.Context(), "unsupported http method "+r.Method)
	}
}

func (j *ActionRunner) runGet(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, request any) ActionResponse,
	request any,
) {
	err := j.fillRequestFromUrlValues(w, r, request, j.router.RouteParams(r))
	if err == nil {
		err = j.fillRequestFromUrlValues(w, r, request, r.URL.Query())
	}

	if err != nil {
		return
	}

	j.runAction(w, r, action, request)
}

func (j *ActionRunner) runPost(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, request any) ActionResponse,
	request any,
) {
	var err error
	err = j.fillRequestFromUrlValues(w, r, request, j.router.RouteParams(r))
	if err == nil {
		if r.Header.Get("Content-Type") == "application/json" {
			err = j.fillRequestFromBody(w, r, request)
		} else {
			err = j.fillRequestFromUrlValues(w, r, request, r.PostForm)
		}
	}

	if err != nil {
		return
	}
	j.runAction(w, r, action, request)
}

func (j *ActionRunner) runPut(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, request any) ActionResponse,
	request any,
) {
	var err error

	err = j.fillRequestFromUrlValues(w, r, request, j.router.RouteParams(r))
	if err != nil {
		return
	}
	err = j.fillRequestFromBody(w, r, request)

	if err != nil {
		return
	}
	j.runAction(w, r, action, request)
}

func (j *ActionRunner) runAction(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, request any) ActionResponse,
	request any,
) {
	if validator, ok := request.(ValidatableStruct); ok {
		validationErr := validator.Validate(r.Context())
		if validationErr != nil {
			j.jsonWriter.Error(w, r, NewValidationErrorResponse(r.Context(), validationErr))
			return
		}
	}

	response := action(r.Context(), request)

	if response.Error != nil {
		j.jsonWriter.Error(w, r, response)
		return
	}
	j.jsonWriter.Success(w, r, response)
}

func (j *ActionRunner) fillRequestFromBody(
	w http.ResponseWriter,
	r *http.Request,
	request any,
) error {
	if request == nil {
		return nil
	}
	if r.Header.Get("Content-Type") != "application/json" {
		return nil
	}
	var err error
	defer r.Body.Close()

	var body []byte

	body, err = io.ReadAll(r.Body)

	if err == nil && body != nil {
		err = json.Unmarshal(body, request)
	}
	if err != nil {
		j.jsonWriter.Error(w, r, NewServerErrorResponse(r.Context(), WrongRequestDecoding, err))
		return err
	}

	return nil
}

func (j *ActionRunner) fillRequestFromUrlValues(
	w http.ResponseWriter,
	r *http.Request,
	request any,
	values url.Values,
) error {
	if request == nil {
		return nil
	}
	err := qs.Unmarshal(request, values.Encode())
	if err != nil {
		resp := j.parseQsError(r.Context(), err)
		j.jsonWriter.Error(w, r, resp)
		return err
	}

	return nil
}

func (j *ActionRunner) parseQsError(ctx context.Context, err error) ActionResponse {
	vErr := NewValidationError("", "Wrong format", InvalidRequest)
	entries := qsErrRegexp.FindStringSubmatch(err.Error())
	if len(entries) > 1 {
		vErr.Field = entries[1]
		if len(entries) > 2 && entries[2] == "strconv.ParseInt" {
			vErr.Err = "Should be integer"
		}
	}
	return NewValidationErrorResponse(ctx, []ValidationError{*vErr})
}
