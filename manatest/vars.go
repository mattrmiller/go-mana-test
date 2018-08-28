// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"fmt"
	"github.com/mattrmiller/go-bedrock/brstrings"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Replace variables in a string.
func ReplaceVars(str string, vars *[]ProjectGlobal) string {

	// Replace global variables
	str = ReplaceGlobalVars(str, vars)

	// Replace random strings
	str = ReplaceRandomString(str)
	str = ReplaceRandomStringLower(str)
	str = ReplaceRandomStringUpper(str)

	// Replace random number
	str = ReplaceRandomNumber(str)

	// Replace cache
	str = ReplaceCache(str)

	return str
}

// Replace global variables in a string.
func ReplaceGlobalVars(str string, vars *[]ProjectGlobal) string {

	// Only replace if we have our context
	if strings.Contains(str, "{{") && len(*vars) != 0 {
		for _, val := range *vars {
			replace := fmt.Sprintf("{{globals.%s}}", val.Key)
			str = strings.Replace(str, replace, val.Value, -1)
		}
	}

	return str
}

// Replace random string.
func ReplaceRandomString(str string) string {

	// Regex
	re, _ := regexp.Compile("{{rand.string.([0-9]*)}}")
	result := re.FindAllStringSubmatch(str, -1)
	if result != nil {
		for _, v := range result {

			// -- Convert to number
			num, _ := strconv.Atoi(v[1])

			// -- Generate random string
			replace := brstrings.RandomAlphaNumString(num)

			// -- Replace
			str = strings.Replace(str, fmt.Sprintf("{{rand.string.%d}}", num), replace, -1)
		}
	}

	return str
}

// Replace random string lower.
func ReplaceRandomStringLower(str string) string {

	// Regex
	re, _ := regexp.Compile("{{rand.string.lower.([0-9]*)}}")
	result := re.FindAllStringSubmatch(str, -1)
	if result != nil {
		for _, v := range result {

			// -- Convert to number
			num, _ := strconv.Atoi(v[1])

			// -- Generate random string
			replace := brstrings.RandomString(num, "abcdefghijklmnopqrstuvwxyz0123456789")

			// -- Replace
			str = strings.Replace(str, fmt.Sprintf("{{rand.string.lower.%d}}", num), replace, -1)
		}
	}
	return str
}

// Replace random string upper.
func ReplaceRandomStringUpper(str string) string {

	// Regex
	re, _ := regexp.Compile("{{rand.string.upper.([0-9]*)}}")
	result := re.FindAllStringSubmatch(str, -1)
	if result != nil {
		for _, v := range result {

			// -- Convert to number
			num, _ := strconv.Atoi(v[1])

			// -- Generate random string
			replace := brstrings.RandomString(num, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

			// -- Replace
			str = strings.Replace(str, fmt.Sprintf("{{rand.string.upper.%d}}", num), replace, -1)
		}
	}
	return str
}

// Replace random number.
func ReplaceRandomNumber(str string) string {

	// Regex
	re, _ := regexp.Compile("{{rand.num.([0-9]*).([0-9]*)}}")
	result := re.FindAllStringSubmatch(str, -1)
	if result != nil {
		for _, v := range result {

			// -- Convert to number
			min, _ := strconv.Atoi(v[1])
			max, _ := strconv.Atoi(v[2])

			// -- Generate random number
			rand.Seed(time.Now().Unix())
			replace := rand.Intn(max-min) + min

			// -- Replace
			str = strings.Replace(str, fmt.Sprintf("{{rand.num.%d.%d}}", min, max), strconv.Itoa(replace), -1)
		}
	}

	return str
}

// Replace cache.
func ReplaceCache(str string) string {

	// Regex
	re, _ := regexp.Compile("{{cache.([^}}]*)}}")
	result := re.FindAllStringSubmatch(str, -1)
	if result != nil {
		for _, v := range result {

			// -- Get cache
			cacheValue := GetCache(v[1])

			// -- Replace
			str = strings.Replace(str, fmt.Sprintf("{{cache.%s}}", v[1]), cacheValue, -1)
		}
	}

	return str
}
