// Package brstrings provides useful string functions.
package brstrings

// Imports
import (
	"strconv"
	"strings"
)

// StringToBool Converts string to bool.
func StringToBool(value string) (bool, error) {

	// Parse float
	ret, err := strconv.ParseBool(strings.ToLower(value))
	if err != nil {
		return false, err
	}

	return ret, err
}

// MustStringToBool Converts string to bool. Panics on error.
func MustStringToBool(value string) bool {

	// Parse float
	ret, err := StringToBool(value)
	if err != nil {
		panic(err)
	}

	return ret
}

// StringToInt Converts string to int.
func StringToInt(value string) (int, error) {

	// Convert to Int
	ret, err := strconv.Atoi(value)
	if err != nil {
		return int(0), err
	}

	return ret, nil
}

// MustStringToInt Converts string to int. Panics on error.
func MustStringToInt(value string) int {

	// Convert to Int
	ret, err := StringToInt(value)
	if err != nil {
		panic(err)
	}

	return ret
}

// StringToInt64 Converts string to int64.
func StringToInt64(value string) (int64, error) {

	// Parse int
	ret, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return int64(0), err
	}

	return ret, nil
}

// MustStringToInt64 Converts string to int64. Panics on error.
func MustStringToInt64(value string) int64 {

	// Parse int
	ret, err := StringToInt64(value)
	if err != nil {
		panic(err)
	}

	return ret
}

// StringToFloat Converts string to float.
func StringToFloat(value string) (float64, error) {

	// Parse float
	ret, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return float64(0), err
	}

	return ret, nil
}

// MustStringToFloat Converts string to float. Panics on error.
func MustStringToFloat(value string) float64 {

	// Parse float
	ret, err := StringToFloat(value)
	if err != nil {
		panic(err)
	}

	return ret
}
