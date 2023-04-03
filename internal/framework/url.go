package framework

import (
	"fmt"
	"net/url"
)

type Url struct {
	originUrl string
	urlObject *url.URL
}

func NewUrl(originUrl string) (*Url, error) {
	urlObject, err := url.Parse(originUrl)
	if err != nil {
		return nil, err
	}

	return &Url{originUrl: originUrl, urlObject: urlObject}, nil
}

func (u *Url) HasBasicAuth() bool {
	return u.urlObject.User != nil
}

func (u *Url) GetBasicAuth() (string, string) {
	if !u.HasBasicAuth() {
		return "", ""
	}
	password, _ := u.urlObject.User.Password()

	return u.urlObject.User.Username(), password
}

func (u *Url) GetUrlWithoutBasicAuth() string {
	result := fmt.Sprintf(
		"%s://%s%s",
		u.urlObject.Scheme,
		u.urlObject.Host,
		u.urlObject.Path,
	)
	if u.urlObject.RawQuery != "" {
		result += "?" + u.urlObject.RawQuery
	}
	if u.urlObject.Fragment != "" {
		result += "#" + u.urlObject.Fragment
	}

	return result
}
