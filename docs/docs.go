package docs

import (
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

// ReadSpec return string with openapi specification
func ReadSpec() string {
	spec, err := ioutil.ReadFile("docs/openapi.json")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't read openapi specification file"))
	}
	return string(spec)
}
