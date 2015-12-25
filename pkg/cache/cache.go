package cache



//return status of chache
type CacheStatus struct {
	Gets        int64
	Hits        int64
	MaxItemSize int
	CurrentSize int
}

//this is a interface which defines some common functions
type Cache interface{
	Set(key interface{}, value interface{})
	Get(key interface{}) (interface{}, bool)
	Delete(key interface{})
	Status()(*CacheStatus)
}