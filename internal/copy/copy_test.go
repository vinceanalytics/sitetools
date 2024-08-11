package copy

import "testing"

func TestCopy(t *testing.T) {
	err := Copy(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
}
