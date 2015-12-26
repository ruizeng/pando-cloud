package cache

import (
	"testing"
)

type simpleStruct struct {
	int
	string
}

type complexStruct struct {
	int
	simpleStruct
}

var getTests = []struct {
	name       string
	keyToAdd   string
	keyToGet   string
	expectedOk bool
}{
	{"string_hit", "myKey", "myKey", true},
	{"string_miss", "myKey", "nonsense", false},
}

func TestSet(t *testing.T) {
	var cache Cache
	cache = NewMemCache(0)
	values := []string{"test1", "test2", "test3"}
	key := "key1"
	for _, v := range values {
		cache.Set(key, v)
		val, ok := cache.Get(key)
		if !ok{
			t.Fatalf("expect key:%v ,value:%v", key, v)
		} else if ok && val != v {
			t.Fatalf("expect value:%v, get value:%v", key, v, val)
		}
		t.Logf("value:%v ", val)
	}
}

func TestGet(t *testing.T) {
	var cache Cache
	cache = NewMemCache(0)
	for _, tt := range getTests {
		cache.Set(tt.keyToAdd, 1234)
		val, ok := cache.Get(tt.keyToGet)

		if ok != tt.expectedOk {
			t.Fatalf("%s: val:%v cache hit = %v; want %v", tt.name, val, ok, !ok)
		} else if ok && val != 1234 {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}

	}
}

func TestDelete(t *testing.T) {
	var cache Cache
	cache = NewMemCache(0)
	cache.Set("myKey", 1234)
	if val, ok := cache.Get("myKey"); !ok {
		t.Fatal("TestRemove returned no match")
	} else if val != 1234 {
		t.Fatalf("TestRemove failed.  Expected %d, got %v", 1234, val)
	}

	cache.Delete("myKey")
	if _, ok := cache.Get("myKey"); ok {
		t.Fatal("TestRemove returned a removed item")
	}
}

func TestStatus(t *testing.T) {

	keys := []string{"1", "2", "3", "4", "5"}

	var gets int64
	var hits int64
	var maxSize int
	var currentSize int
	maxSize = 20
	var cache Cache
	cache = NewMemCache(maxSize)
	//keys := []string{"1", "2", "3", "4", "5"}
	for _, key := range keys {
		cache.Set(key, 1234)
		currentSize++
	}

	newKeys := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	for _, newKey := range newKeys {
		_, ok := cache.Get(newKey)
		if ok == true {
			hits++
		}
		gets++
	}
	t.Logf("gets:%v, hits:%v, maxSize:%v, currentSize:%v", gets, hits, maxSize, currentSize)
	status := cache.Status()
	if status.CurrentSize != currentSize || status.MaxItemSize != maxSize ||
		status.Gets != gets || status.Hits != hits {
		t.Fatalf("get status maxSize:%v, currentSize:%v, nget:%v, nhit:%v",
			status.MaxItemSize, status.CurrentSize, status.Gets, status.Hits)
	}

}

func TestLRU(t *testing.T) {
	keys := []string{"1", "2", "3", "4", "2", "1", "3", "5", "6", "5", "6"}
	maxSize := 3
	var cache Cache
	cache = NewMemCache(maxSize)
	for i, key := range keys {
		cache.Set(key, 1234)
		if i == 3 {
			status := cache.Status()
			if status.CurrentSize != maxSize {
				t.Fatalf("expected maxSize %v,currentSize:%v", maxSize, status.CurrentSize)
			}
			_, ok1 := cache.Get("2")
			_, ok2 := cache.Get("3")
			_, ok3 := cache.Get("4")
			if !(ok1 && ok2 && ok3) {
				t.Fatalf("expected remains key 2:%v,3:%v, 4:%v", ok1, ok2, ok3)
			}
		}
	}

	status := cache.Status()
	if status.CurrentSize != maxSize {
		t.Fatalf("expected maxSize %v,currentSize:%v", maxSize, status.CurrentSize)
	}

	_, ok1 := cache.Get("3")
	_, ok2 := cache.Get("5")
	_, ok3 := cache.Get("6")
	if !(ok1 && ok2 && ok3) {
		t.Fatalf("expected remains key 3:%v,5:%v, 6:%v", ok1, ok2, ok3)
	}
}
