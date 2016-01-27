package muts

import "fmt"

// Tos returns the string representation of the value
func Tos(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
