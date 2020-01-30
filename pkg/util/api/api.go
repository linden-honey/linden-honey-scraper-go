package api

import (
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

// ReadSpec return specification data
func ReadSpec(path string) string {
	spec, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't read openapi specification file"))
	}
	return string(spec)
}
