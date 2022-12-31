package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type VoiceParam struct {
	Lang     string
	Text     string
	Instance string
}

type VoiceData struct {
	Audio []byte
}

func NewVoiceParam(lang, text, instance string) (*VoiceParam, error) {
	if strings.TrimSpace(lang) == "" {
		return nil, fmt.Errorf("%w: language must not be a empty string", ErrInvalidArgs)
	}

	return &VoiceParam{
		Lang:     lang,
		Text:     text,
		Instance: instance,
	}, nil
}

func (param *VoiceParam) ToURL() *url.URL {
	if strings.TrimSpace(param.Instance) == "" {
		param.Instance = defaultInstance
	}

	//nolint:exhaustivestruct,exhaustruct
	return &url.URL{
		Scheme: "https",
		Host:   param.Instance,
		Path:   path.Join("api", "v1", "audio", param.Lang, param.Text),
	}
}

func (param *VoiceParam) ToHTTPRequest() (*http.Request, error) {
	//nolint:noctx
	req, err := http.NewRequest(http.MethodGet, param.ToURL().String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequest, err)
	}

	return req, nil
}

func FetchVoice(param *VoiceParam) (*VoiceData, error) {
	req, err := param.ToHTTPRequest()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequest, err)
	}

	body, err := httpGet(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrResponse, err)
	}

	data, err := parseVoiceData(body)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrResponse, err)
	}

	return data, nil
}

func parseVoiceData(body io.ReadCloser) (*VoiceData, error) {
	var data VoiceData
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrResponse, err)
	}

	defer body.Close()

	return &data, nil
}
