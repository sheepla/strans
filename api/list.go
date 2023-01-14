package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/tidwall/gjson"
)

type Lang struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func ListSourceLangs(instance string) (*[]Lang, error) {
	u := newSourceLangsURL(instance)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: (url=%s) %s", ErrRequest, u.String(), err)
	}

	body, err := httpGet(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get source languages: %w: %s", ErrResponse, err)
	}

	defer body.Close()

	langs, err := parseLangsResult(body)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParse, err)
	}

	return langs, nil
}

func newSourceLangsURL(instance string) *url.URL {
	if strings.TrimSpace(instance) == "" {
		instance = defaultInstance
	}

	u := &url.URL{
		Scheme: "https",
		Host:   instance,
		Path:   path.Join("api", "v1", "languages", "source"),
	}

	return u
}

func ListTargetLangs(instance string) (*[]Lang, error) {
	u := newTargetLangsURL(instance)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: (url=%s) %s", ErrRequest, u.String(), err)
	}

	body, err := httpGet(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get target languages: %w: %s", ErrResponse, err)
	}

	defer body.Close()

	langs, err := parseLangsResult(body)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParse, err)
	}

	return langs, nil
}

func newTargetLangsURL(instance string) *url.URL {
	if strings.TrimSpace(instance) == "" {
		instance = defaultInstance
	}

	u := &url.URL{
		Scheme: "https",
		Host:   instance,
		Path:   path.Join("api", "v1", "languages", "target"),
	}

	return u
}

func parseLangsResult(r io.Reader) (*[]Lang, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var langs []Lang
	gjson.GetBytes(raw, "languages").ForEach(func(key, value gjson.Result) bool {
		langs = append(langs, Lang{
			Code: value.Get("code").String(),
			Name: value.Get("name").String(),
		})

		return true
	})

	return &langs, nil

}
