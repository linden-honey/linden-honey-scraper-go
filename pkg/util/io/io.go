package io

import (
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

// MustReadContent return file's data
func MustReadContent(path string) string {
	spec, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can't read file"))
	}
	return string(spec)
}
