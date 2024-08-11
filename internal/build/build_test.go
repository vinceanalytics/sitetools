package build

import (
	"bytes"
	"os"
	"testing"
)

func TestJSONPre(t *testing.T) {
	var page Page
	f, err := os.ReadFile("testdata/page.md")
	if err != nil {
		t.Fatal(err)
	}

	data, err := jsonPre(&page, f)
	if err != nil {
		t.Fatal(err)
	}
	front := []byte("A basic markdown page with json front matter")
	if !bytes.Equal(data, front) {
		t.Fatal("wrong front matter")
	}
}
