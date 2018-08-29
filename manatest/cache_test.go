// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"github.com/mattrmiller/go-bedrock/brtesting"
	"testing"
)

// Test ClearCache.
func TestClearCache(tst *testing.T) {

	// Clear cache
	ClearCache()

	// Set Some cache
	SetCache("one", "1")
	brtesting.AssertEqual(tst, GetCache("one"), "1", "ClearCache failed for 'one'")

	// Test non existent
	ClearCache()
	brtesting.AssertEqual(tst, GetCache("one"), "", "ClearCache failed for 'one'")
}

// Test GetCacheKeys.
func TestGetCacheKeys(tst *testing.T) {

	// Clear cache
	ClearCache()

	// Set Some cache
	SetCache("one", "1")
	SetCache("two", "1")
	keys := GetCacheKeys()
	brtesting.AssertEqual(tst, keys[0], "one", "GetCacheKeys failed for 'one'")
	brtesting.AssertEqual(tst, keys[1], "two", "GetCacheKeys failed for 'one'")
}

// Test GetCache.
func TestGetCache(tst *testing.T) {

	// Clear cache
	ClearCache()

	// Set Some cache
	SetCache("two", "2")
	brtesting.AssertEqual(tst, GetCache("two"), "2", "GetCache failed for 'one'")

	// Test non existent
	brtesting.AssertEqual(tst, GetCache("none"), "", "GetCache failed for 'none'")
}
