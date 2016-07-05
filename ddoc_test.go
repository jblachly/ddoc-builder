package ddoc

import (
	"fmt"
	"testing"
)

func TestBuild(t *testing.T) {
	b, err := Build("design", "couchdb")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

}
