// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"fmt"
	"github.com/mattrmiller/go-mana-test/console"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
	"reflect"
	"strings"
)

// Cache
const (
	CacheBodyJSON = "response.body.json"
)

// Global cache
var cache map[string]string

// ValidateCacheValue Validates cache value.
func ValidateCacheValue(value *string) bool {

	// Body json
	if strings.HasPrefix(*value, CacheBodyJSON) {
		return true
	}

	return false
}

// SaveCacheFromResponse Saves cache from response.
func SaveCacheFromResponse(caches *[]TestCache, response *resty.Response) error {

	// Each check
	for _, cache := range *caches {

		// -- Verbose
		console.PrintVerbose(fmt.Sprintf("\t\t\t- %s", cache.Name))

		// -- Check body json
		if strings.HasPrefix(cache.Value, CheckResBodyJSON) {
			err := cacheBodyJSON(&cache, response)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ClearCache Clears cache
func ClearCache() {
	cache = make(map[string]string)
}

// GetCacheKeys Gets a value from inside of cache.
func GetCacheKeys() []string {
	keys := reflect.ValueOf(cache).MapKeys()
	ret := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		ret[i] = keys[i].String()
	}

	return ret
}

// GetCache Gets a value from inside of cache.
func GetCache(key string) string {
	return cache[key]
}

// SetCache Sets a value inside of cache.
func SetCache(key string, value string) {

	// Make cache
	if cache == nil {
		ClearCache()
	}

	// Set cache
	cache[key] = value
}

// cacheBodyJSON Caches the body json for a response.
func cacheBodyJSON(cache *TestCache, response *resty.Response) error {

	// First make sure that response type was json
	if !strings.Contains(response.Header().Get(HeaderContentType), ContentTypeJSON) {
		return fmt.Errorf("Response '%s' was not a proper content type '%s'", cache.Name, ContentTypeJSON)
	}

	// Scrape the prefix off the selector
	jsonSel := strings.TrimPrefix(cache.Value, fmt.Sprintf("%s.", CheckResBodyJSON))

	// Get the json
	jsonValue := gjson.Get(string(response.Body()), jsonSel)

	// Cache
	if !jsonValue.Exists() {
		return fmt.Errorf("Cache '%s' was null", cache.Value)
	}
	SetCache(cache.Name, jsonValue.Str)

	return nil
}
