package common

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"log"
	"strings"
)

// ErrBadTransform indicates an uneven transformation list
var ErrBadTransform = errors.New("Path transformation list not even. " +
	"Transformations must be specified as pairs.")

// PathTransformer is an instance of the path transformer
type PathTransformer struct {
	pathReplacer *strings.Replacer
}

// NewTransformer returns a new path transformer
func NewTransformer(str string) (*PathTransformer, error) {
	// build replacer instance
	t := new(PathTransformer)

	// check string
	if str != "" {
		// split strings
		split := strings.Split(str, ",")

		// transformations must be even pairs
		if (len(split) % 2) != 0 {
			return nil, ErrBadTransform
		}

		// build replacer
		t.pathReplacer = strings.NewReplacer(split...)
	}

	// all good
	return t, nil
}

// Transform performs path transformation as requested
func (t *PathTransformer) Transform(kvps api.KVPairs) {
	// check replacer - return immediately if not valid
	if t.pathReplacer == nil {
		// do nothing
		return
	}

	// loop through keys
	for _, kv := range kvps {
		// split path and key with strings because
		// the path package will trim a trailing / which
		// breaks empty folders present in the kvp store
		split := strings.Split(kv.Key, ConsulSeparator)
		// get and check length ... only continue if we actually
		// have a path we may want to transform
		if length := len(split); length > 1 {
			// isolate and replace path
			rpath := t.pathReplacer.Replace(strings.Join(split[:length-1], ConsulSeparator))
			// join replaced path with key
			newKey := strings.Join([]string{rpath, split[length-1]}, ConsulSeparator)
			// check keys
			if kv.Key != newKey {
				// log change
				log.Printf("[Transform] %s -> %s", kv.Key, newKey)
				// update key
				kv.Key = newKey
			}
		}
	}
}
