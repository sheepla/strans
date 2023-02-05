package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

var (
	ErrInvalidArgs = errors.New("invalid argument(s)")
	ErrRequest     = errors.New("an error occurred on creating request")
	ErrHTTP        = errors.New("an error occurred on executing HTTP method")
	ErrResponse    = errors.New("an error occurred on processing response")
	ErrParse       = errors.New("an error occurred on parsing result")
)

func Translate(param *TranslateParam) (*Result, error) {
	req, err := param.ToHTTPRequest()
	if err != nil {
		return nil, err
	}

	body, err := httpGet(req)
	if err != nil {
		return nil, err
	}

	result, err := parseResult(body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type TranslateParam struct {
	SourceLang string
	TargetLang string
	Text       string
	Instance   string
}

func NewTranslateParam(source, target, text, instance string) (*TranslateParam, error) {
	if strings.TrimSpace(source) == "" {
		return nil, fmt.Errorf("%w: source must not be empty string", ErrInvalidArgs)
	}

	if strings.TrimSpace(target) == "" {
		return nil, fmt.Errorf("%w: target must not be empty string", ErrInvalidArgs)
	}

	// if strings.TrimSpace(text) == "" {
	// 	return nil, fmt.Errorf("%w: text must not be empty string", ErrInvalidArgs)
	// }

	return &TranslateParam{
		SourceLang: source,
		TargetLang: target,
		Text:       text,
		Instance:   instance,
	}, nil
}

func (param *TranslateParam) ToURL() *url.URL {
	if strings.TrimSpace(param.Instance) == "" {
		param.Instance = defaultInstance
	}

	//nolint:exhaustivestruct,exhaustruct
	return &url.URL{
		Scheme: "https",
		Host:   param.Instance,
		Path: path.Join(
			"api",
			"v1",
			url.PathEscape(param.SourceLang),
			url.PathEscape(param.TargetLang),
			url.PathEscape(param.Text),
		),
	}
}

func (param *TranslateParam) ToHTTPRequest() (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		param.ToURL().String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequest, err)
	}

	return req, nil
}

type Result struct {
	Text string `json:"translation"`
	//	Info struct {
	//		Definitions       []interface{} `json:"definitions"`
	//		Examples          []interface{} `json:"examples"`
	//		ExtraTranslations []struct {
	//			List []struct {
	//				Frequency int64    `json:"frequency"`
	//				Meanings  []string `json:"meanings"`
	//				Word      string   `json:"word"`
	//			} `json:"list"`
	//			Type string `json:"type"`
	//		} `json:"extraTranslations"`
	//		Pronunciation struct {
	//			Translation string `json:"translation"`
	//		} `json:"pronunciation"`
	//		Similar []interface{} `json:"similar"`
	//	} `json:"info"`
}

func parseResult(body io.ReadCloser) (*Result, error) {
	var result Result
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrResponse, err)
	}

	defer body.Close()

	// unescape url
	if escaped, err := url.PathUnescape(result.Text); err == nil {
		result.Text = escaped
	}

	return &result, nil
}
