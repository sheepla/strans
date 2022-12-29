//nolint:godot,varnamelen
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultInstance = "simplytranslate.org"
	timeout         = 10
)

var (
	ErrArgIsEmpty  = errors.New("argument is empty")
	ErrParseEngine = errors.New("failed to parse engine string")
	ErrAPI         = errors.New("an error occurred when calling the API")
)

type Param struct {
	SourceLang string
	TargetLang string
	Text       string
	Engine     Engine
}

func NewParam(source, target, text string, engine Engine) (*Param, error) {
	if strings.TrimSpace(source) == "" {
		return nil, fmt.Errorf("%w (source language)", ErrArgIsEmpty)
	}

	if strings.TrimSpace(target) == "" {
		return nil, fmt.Errorf("%w (target language)", ErrArgIsEmpty)
	}

	if strings.TrimSpace(target) == "" {
		return nil, fmt.Errorf("%w (text)", ErrArgIsEmpty)
	}

	param := &Param{
		SourceLang: source,
		TargetLang: target,
		Text:       text,
		Engine:     engine,
	}

	return param, nil
}

type Engine int

const (
	EngineDefault Engine = iota
	EngineGoogle
	EngineLibre
	EngineDeepL
	// EngineICIBA
	// EngineReverso
)

func (e Engine) String() string {
	switch e {
	case EngineGoogle:
		return "google"
	case EngineLibre:
		return "libre"
	case EngineDeepL:
		return "deepl"
	case EngineDefault:
		return ""
	default:
		return ""
	}
}

func ParseEngineString(s string) (Engine, error) {
	lower := strings.ToLower(strings.TrimSpace(s))
	switch lower {
	case "google":
		return EngineGoogle, nil
	case "libre":
		return EngineDeepL, nil
	case "deepl":
		return EngineDeepL, nil
	case "default":
		return EngineDefault, nil
	default:
		return EngineDefault, fmt.Errorf("%w (%s)", ErrParseEngine, s)
	}
}

//nolint:tagliatelle
type Result struct {
	SourceLang string `json:"source_language"`
	Text       string `json:"translated-text"`
}

func Translate(param *Param, instance string) (*Result, error) {
	req, err := newRequest(param, instance)
	if err != nil {
		return nil, err
	}

	body, err := fetch(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrAPI, err)
	}

	result, err := parseResult(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the result: %w", err)
	}

	return result, nil
}

func parseResult(body io.ReadCloser) (*Result, error) {
	var result Result

	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response as JSON: %w", err)
	}

	defer body.Close()

	return &result, nil
}

func fetch(req *http.Request) (io.ReadCloser, error) {
	//nolint:exhaustivestruct,exhaustruct
	c := http.Client{
		Timeout: timeout * time.Second,
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send a request: %w", err)
	}

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		//nolint:goerr113
		return nil, fmt.Errorf("HTTP status error: %s", resp.Status)
	}

	return resp.Body, nil
}

func newURL(param *Param, instance string) *url.URL {
	if strings.TrimSpace(instance) == "" {
		instance = defaultInstance
	}

	//nolint:exhaustivestruct,exhaustruct,varnamelen
	u := &url.URL{
		Scheme: "https",
		Host:   instance,
		Path:   "api/translate",
	}

	q := u.Query()
	q.Add("from", param.SourceLang)
	q.Add("to", param.TargetLang)
	q.Add("text", param.Text)

	if param.Engine != EngineDefault {
		q.Add("engine", param.Engine.String())
	}

	u.RawQuery = q.Encode()

	return u
}

func newRequest(param *Param, instance string) (*http.Request, error) {
	//nolint:varnamelen
	u := newURL(param, instance)

	//nolint:noctx
	req, err := http.NewRequest(
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	return req, nil
}
