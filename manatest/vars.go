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

	// Replace random string
	str = ReplaceRandomString(str)

	// Replace random number
	str = ReplaceRandomNumber(str)

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
	result := re.FindStringSubmatch(str)
	for _, v := range result {

		// -- Convert to number
		num, _ := strconv.Atoi(v)

		// -- Generate random string
		replace := brstrings.RandomAlphaNumString(num)

		// -- Replace
		str = strings.Replace(str, fmt.Sprintf("{{rand.string.%d}}", num), replace, -1)
	}

	return str
}

// Replace random number.
func ReplaceRandomNumber(str string) string {

	// Regex
	re, _ := regexp.Compile("{{rand.num.([0-9]*).([0-9]*)}}")
	result := re.FindAllStringSubmatch(str, -1)
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

	return str
}
