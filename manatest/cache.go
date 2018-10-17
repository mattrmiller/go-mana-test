// Package manatest provides internal workings for go-mana-test.
package manatest

// Imports
import (
	"fmt"
	"github.com/mattrmiller/go-mana-test/http"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
	"reflect"
	"strings"
)

// Cache.
const (
	CacheBodyJSON = "response.body.json"
)

// Global cache.
var cache map[string]string

// ValidateCacheValue Validates cache value.
func ValidateCacheValue(value *string) bool {

	// Body json
	return strings.HasPrefix(*value, CacheBodyJSON)
}

// SaveCacheFromResponse Saves cache from response.
func SaveCacheFromResponse(caches *[]TestCache, response *resty.Response) error {

	// Each check
	for _, chce := range *caches {

		// Check body json
		if strings.HasPrefix(chce.Value, CacheBodyJSON) {
			err := cacheBodyJSON(&chce, response)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ClearCache Clears cache.
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
	if !strings.Contains(response.Header().Get(http.HeaderContentType), http.ContentTypeJSON) {
		return fmt.Errorf("Response '%s' was not a proper content type '%s'", cache.Name, http.ContentTypeJSON)
	}

	// Scrape the prefix off the selector
	jsonSel := strings.TrimPrefix(cache.Value, fmt.Sprintf("%s.", CheckResBodyJSON))

	// Get the json
	jsonValue := gjson.Get(string(response.Body()), jsonSel)

	// Cache
	if !jsonValue.Exists() {
		return fmt.Errorf("Cache '%s' was null", cache.Value)
	}
	SetCache(cache.Name, jsonValue.String())

	return nil
}
