package cache

//this is a interface which defines some common functions
type Cache interface{
	Set(key interface{}, value interface{})
	Get(key interface{}) (interface{}, bool)
	Delete(key interface{})
	Status()(interface{})
}