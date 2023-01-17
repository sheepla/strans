// nolint
package api

import "testing"

var p1 = &TranslateParam{
	SourceLang: "en",
	TargetLang: "ja",
	Text:       "The quick brown fox jumps over the lazy dog. Hello/World",
	Instance:   "",
}

func TestTranslate(t *testing.T) {
	result, err := Translate(p1)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}

func TestParamToURL(t *testing.T) {
	u := p1.ToURL()
	have := u.String()
	want := `https://lingva.ml/api/v1/en/ja/The%2520quick%2520brown%2520fox%2520jumps%2520over%2520the%2520lazy%2520dog.%2520Hello%252FWorld`
	if have != want {
		t.Fatal("have:", have, " want:", want)
	}

	t.Log(have)
}
