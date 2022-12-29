//nolint
package api

import "testing"

var p1 = &Param{
	SourceLang: "en",
	TargetLang: "ja",
	Text:       "The quick brown fox jumps over the lazy dog.",
}

var instance1 = "translate.tiekoetter.com"

func TestTranslate(t *testing.T) {
	result, err := Translate(p1, "")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)

	result, err = Translate(p1, instance1)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}

func TestNewRequest(t *testing.T) {
	req, err := newRequest(p1, "")
	if err != nil {
		t.Errorf(err.Error())
	}

	t.Log(req)
}

func TestNewURL(t *testing.T) {
	u := newURL(p1, "")
	t.Log(u.String())
}

func TestParseEngineString(t *testing.T) {
	e, err := ParseEngineString("google")
	t.Log(e)
	if err != nil {
		t.Fatal(err)
	}

	e, err = ParseEngineString("UNKNOWN ENGINE")
	if err == nil {
		t.Fatal("err is nil")
	}
}
