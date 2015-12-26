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
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string)
	Status()(*CacheStatus)
}