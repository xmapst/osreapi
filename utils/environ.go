package utils

import "sort"

// Slice is a helper function that converts a map of environment
// variables to a slice of string values in key=value format.
func Slice(env map[string]string) []string {
	var s []string
	for k, v := range env {
		s = append(s, k+"="+v)
	}
	sort.Strings(s)
	return s
}

// Combine is a helper function combines one or more maps of
// environment variables into a single map.
func Combine(env ...map[string]string) map[string]string {
	c := map[string]string{}
	for _, e := range env {
		for k, v := range e {
			c[k] = v
		}
	}
	return c
}
