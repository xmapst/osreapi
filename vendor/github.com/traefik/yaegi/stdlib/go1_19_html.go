// Code generated by 'yaegi extract html'. DO NOT EDIT.

//go:build go1.19 && !go1.20
// +build go1.19,!go1.20

package stdlib

import (
	"html"
	"reflect"
)

func init() {
	Symbols["html/html"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"EscapeString":   reflect.ValueOf(html.EscapeString),
		"UnescapeString": reflect.ValueOf(html.UnescapeString),
	}
}
