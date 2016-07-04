package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type DesignDoc struct {
	ID  string `json:"_id"`
	Rev string `json:"_rev,omitempty"`

	Language string                 `json:"language"`
	Views    map[string]interface{} `json:"views"`
}

// type View is used to fill the design doc template
// with javascript files
// it is not marshaled to or from JSON itself
type View struct {
	Name           string
	MapFunction    string
	ReduceFunction string
}

var views []View

// BuildDesignDocument takes a design document name (do not include _design/)
// and a file path where javascript views can be found; these views should be
// map.js (TODO: reduce.js) files under subdirectories named after the view
//
// Example:
// path/views/customers/map.js
//
// From these, it composes a CouchDB design document having the map and reduce
// functions as named views
func BuildDesignDocument(name, path string) error {

	if strings.ContainsRune(name, '/') {
		return errors.New("Do not include / in design document name")
	} else if name[0:7] == "_design" {
		return errors.New("Do not begin design document name with _design")
	}

	ddocID := "_design/" + name
	//ddoc := &DesignDoc{ID: ddocID}
	_ = &DesignDoc{ID: ddocID}

	// Glob files
	mapFiles := filepath.Join(path, "views/*/", "map.js")

	matches, err := filepath.Glob(mapFiles)
	if err != nil {
		// only if malformed pattern
		log.Fatalln(err)
	}
	if len(matches) < 1 {
		log.Fatalln("No javascript views found in " + path)
	}

	// Build views
	for i := 0; i < len(matches); i++ {
		v := View{}
		// take the last directory from the path by removing filename (.Dir)
		// and then taking the last element (.Base)
		v.Name = filepath.Base(filepath.Dir(matches[i]))

		// read map.js
		filebuf, err := ioutil.ReadFile(matches[i])
		if err != nil {
			return err
		}
		// TODO: replace \n

		v.MapFunction = string(filebuf)
		v.ReduceFunction = "_sum"

		// Append to 'views' slice
		views = append(views, v)
	}

	// Fill template
	// TODO use passed path
	t, err := template.ParseFiles("couchdb/ddoc.tmpl")
	if err != nil {
		log.Fatalln(err)
	}
	// print results
	t.Execute(os.Stdout, views)

	return nil
}

func main() {
	_ = BuildDesignDocument("design", "couchdb")
}
