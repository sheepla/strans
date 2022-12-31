//nolint
package trans

import "testing"

var p1 = &TranslateParam{
	SourceLang: "en",
	TargetLang: "ja",
	Text:       "The quick brown fox jumps over the lazy dog.",
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
	want := `https://lingva.ml/api/v1/en/ja/The%20quick%20brown%20fox%20jumps%20over%20the%20lazy%20dog.`
	if have != want {
		t.Fatal("have:", have, " want:", want)
	}

	t.Log(have)
}
