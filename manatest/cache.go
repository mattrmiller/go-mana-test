// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"errors"
	"fmt"
	"github.com/mattrmiller/go-mana-test/console"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
	"reflect"
	"strings"
)

// Cache
const (
	CACHE_BODY_JSON = "response.body.json"
)

// Global cache
var cache map[string]string

// Validates cache value.
func ValidateCacheValue(value *string) bool {

	// Body json
	if strings.HasPrefix(*value, CACHE_BODY_JSON) {
		return true
	}

	return false
}

// Saves cache from response.
func SaveCacheFromResponse(caches *[]TestCache, response *resty.Response) error {

	// Each check
	for _, cache := range *caches {

		// -- Verbose
		console.PrintVerbose(fmt.Sprintf("\t\t\t- %s", cache.Name))

		// -- Check body json
		if strings.HasPrefix(cache.Value, CHECK_RES_BODY_JSON) {
			err := cacheBodyJson(&cache, response)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Gets a value from inside of cache.
func GetCacheKeys() []string {
	keys := reflect.ValueOf(cache).MapKeys()
	ret := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		ret[i] = keys[i].String()
	}

	return ret
}

// Gets a value from inside of cache.
func GetCache(key string) string {
	return cache[key]
}

// Sets a value inside of cache.
func SetCache(key string, value string) {

	// Make cache
	if cache == nil {
		cache = make(map[string]string)
	}

	// Set cache
	cache[key] = value
}

// Caches the body json for a response.
func cacheBodyJson(cache *TestCache, response *resty.Response) error {

	// First make sure that response type was json
	if !strings.Contains(response.Header().Get(HEADER_CONTENT_TYPE), CONTENT_TYPE_JSON) {
		return errors.New(fmt.Sprintf("Response '%s' was not a proper content type '%s'", cache.Name, CONTENT_TYPE_JSON))
	}

	// Scrape the prefix off the selector
	jsonSel := strings.TrimPrefix(cache.Value, fmt.Sprintf("%s.", CHECK_RES_BODY_JSON))

	// Get the json
	jsonValue := gjson.Get(string(response.Body()), jsonSel)

	// Cache
	if !jsonValue.Exists() {
		return errors.New(fmt.Sprintf("Cache '%s' was null", cache.Value))
	}
	SetCache(cache.Name, jsonValue.Str)

	return nil
}
