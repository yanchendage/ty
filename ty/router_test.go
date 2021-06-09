package ty

import (
	"reflect"
	"testing"
)

func TestParsePattern(t *testing.T) {


	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}