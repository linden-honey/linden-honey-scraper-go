package config

import (
	"fmt"
	"reflect"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"

	"github.com/caarlos0/env/v6"
)

func Parsers() map[reflect.Type]env.ParserFunc {
	return map[reflect.Type]env.ParserFunc{
		reflect.TypeOf((*encoding.Encoding)(nil)).Elem(): func(v string) (interface{}, error) {
			e, err := htmlindex.Get(v)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve encoding by name=%s: %w", v, err)
			}

			return e, nil
		},
	}
}
