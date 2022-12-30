package trans

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	defaultInstance = "lingva.ml"
	timeout         = 10 * time.Second
)

var (
	ErrInvalidArgs = errors.New("invalid argument(s)")
	ErrRequest     = errors.New("an error occurred on creating request")
	ErrHTTP        = errors.New("an error occurred on executing HTTP method")
	ErrResponse    = errors.New("an error occurred on processing response")
	ErrAPI         = errors.New("an error occurred on calling API")
)

func Translate(param *Param) (*Result, error) {
	body, err := httpGet(param)
	if err != nil {
		return nil, err
	}

	result, err := parseResult(body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type Param struct {
	SourceLang string
	TargetLang string
	Text       string
	Instance   string
}

func NewParam(source, target, text, instance string) (*Param, error) {
	if strings.TrimSpace(source) == "" {
		return nil, fmt.Errorf("%w: source must not be empty string", ErrInvalidArgs)
	}

	if strings.TrimSpace(target) == "" {
		return nil, fmt.Errorf("%w: target must not be empty string", ErrInvalidArgs)
	}

	// if strings.TrimSpace(text) == "" {
	// 	return nil, fmt.Errorf("%w: text must not be empty string", ErrInvalidArgs)
	// }

	return &Param{
		SourceLang: source,
		TargetLang: target,
		Text:       text,
		Instance:   instance,
	}, nil
}

func (param *Param) ToURL() *url.URL {
	if strings.TrimSpace(param.Instance) == "" {
		param.Instance = defaultInstance
	}

	//nolint:exhaustivestruct,exhaustruct
	return &url.URL{
		Scheme: "https",
		Host:   param.Instance,
		Path:   path.Join("api", "v1", param.SourceLang, param.TargetLang, param.Text),
	}
}

func httpGet(param *Param) (io.ReadCloser, error) {
	//nolint:noctx
	req, err := http.NewRequest(
		http.MethodGet,
		param.ToURL().String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequest, err)
	}

	//nolint:exhaustivestruct,exhaustruct
	cl := &http.Client{
		Timeout: timeout,
	}

	resp, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return nil, fmt.Errorf("%w: %s", ErrHTTP, resp.Status)
	}

	return resp.Body, nil
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

	return &result, nil
}
