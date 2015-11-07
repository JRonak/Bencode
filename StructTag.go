package Bencode

import (
	"errors"
	"reflect"
	"strings"
)

const (
	key   = "Bencode"
	empty = "omitempty"
)

var (
	noTag error
)

func init() {
	noTag = errors.New("No tags")
}

func getTag(r reflect.StructField) string {
	tag := r.Tag.Get(key)
	if tag == "" {
		return r.Name
	}
	if strings.Contains(tag, ",") {
		s := tag[:strings.Index(tag, ",")]
		if s == "" {
			return r.Name
		}
		return s
	} else {
		if tag == empty {
			return r.Name
		} else {
			return tag
		}
	}
}

func isOmitEmpty(r reflect.StructField) bool {
	tag := r.Tag.Get(key)
	if tag == "" {
		return false
	} else if !strings.Contains(tag, ",") {
		if tag == empty {
			return true
		} else {
			return false
		}
	} else {
		e := tag[strings.Index(tag, ",")+1:]
		if empty == e {
			return true
		}
		return false
	}

}

func isIgnore(r reflect.StructField) bool {
	tag := r.Tag.Get(key)
	if tag == "-" {
		return true
	} else {
		return false
	}
}
