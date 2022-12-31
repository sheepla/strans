//nolint
package api

import "testing"

var vp1 = &VoiceParam{
	Lang:     "ja",
	Text:     "こんにちは世界",
	Instance: "",
}

func TestVoiceParamToURL(t *testing.T) {
	have := vp1.ToURL().String()
	want := `https://lingva.ml/api/v1/audio/ja/%E3%81%93%E3%82%93%E3%81%AB%E3%81%A1%E3%81%AF%E4%B8%96%E7%95%8C`
	if have != want {
		t.Fatal("have:", have, "want:", want)
	}
}
