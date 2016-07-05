package ddoc

import (
	"bytes"
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

// views is passed to the template
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
func Build(name, path string) ([]byte, error) {

	// Error checking
	if strings.ContainsRune(name, '/') {
		return nil, errors.New("Do not include / in design document name")
	} else if strings.Contains(name, "_design") {
		return nil, errors.New("Do not begin design document name with _design")
	}

	ddocTemplate := filepath.Join(path, "ddoc.tmpl")
	if _, err := os.Stat(ddocTemplate); os.IsNotExist(err) {
		// path/ddoc.tmpl does not exist
		return nil, err
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
			return nil, err
		}

		mapFunc := string(filebuf)
		mapFunc = strings.Replace(mapFunc, "\n", "\\n", -1) // replace newline with literal \n
		mapFunc = strings.Replace(mapFunc, "\r", "", -1)    // replace carriage ret with nothing

		v.MapFunction = string(mapFunc)
		v.ReduceFunction = "_sum"

		// Append to 'views' slice
		views = append(views, v)
	}

	// Fill template
	t, err := template.ParseFiles(ddocTemplate)
	if err != nil {
		log.Fatalln(err)
	}
	// print results
	buf := make([]byte, 0, 8192) // byte slice with capacity of 8k
	b := bytes.NewBuffer(buf)
	t.Execute(b, views)

	return b.Bytes(), nil
}
