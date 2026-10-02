package up

import "testing"

func Test_toJSON_NilSliceIsEmptyList(t *testing.T) {
	var none []AccountResource
	if got := toJSON(none); got != "[]" {
		t.Errorf("toJSON(nil slice) = %q, want []", got)
	}
	if got := toJSON([]AccountResource{}); got != "[]" {
		t.Errorf("toJSON(empty slice) = %q, want []", got)
	}
	if got := toJSON(map[string]string{"a": "b"}); got == "[]" || got == "null" {
		t.Errorf("toJSON(map) = %q", got)
	}
}
