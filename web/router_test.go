package web

import (
	"reflect"
	"testing"
	"ty"
)

func TestParsePattern(t *testing.T) {


	ok := reflect.DeepEqual(ty.parsePattern("/p/:name"), []string{"p", ":name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}