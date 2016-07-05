package ddoc

import (
	"fmt"
	"testing"
)

func TestBuild(t *testing.T) {
	s, err := Build("design", "couchdb")
	if err != nil {
		panic(err)
	}

	fmt.Println(s)

}
