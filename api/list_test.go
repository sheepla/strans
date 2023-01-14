package api

import "testing"

func TestListTargetLangs(t *testing.T) {
	langs, err := ListSourceLangs("")
	if err != nil {
		t.Fatal(err)
	}

	if langs == nil {
		t.Fatal("langs is nil")
	}

	for _, v := range *langs {
		t.Log(v.Code, " -> ", v.Name)
	}
}

func TestListSourceLangs(t *testing.T) {
	langs, err := ListSourceLangs("")
	if err != nil {
		t.Fatal(err)
	}

	if langs == nil {
		t.Fatal("langs is nil")
	}

	for _, v := range *langs {
		t.Log(v.Code, " -> ", v.Name)
	}
}

func TestNewSourceLangsURL(t *testing.T) {
	have := newSourceLangsURL("").String()
	want := "https://lingva.ml/api/v1/languages/source"
	if have != want {
		t.Fatal("have=", have, "want=", want)
	}
}

func TestNewTargetLangsURL(t *testing.T) {
	have := newTargetLangsURL("").String()
	want := "https://lingva.ml/api/v1/languages/target"
	if have != want {
		t.Fatal("have=", have, "want=", want)
	}
}
