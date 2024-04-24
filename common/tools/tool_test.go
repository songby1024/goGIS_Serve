package tools

import "testing"

func TestParseIntSliceFromString(t *testing.T) {
	str, err := ParseIntSliceFromString("{1,3,5}")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("data:", str)
}
