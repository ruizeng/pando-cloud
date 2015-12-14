package main

import (
	"fmt"
	"reflect"
)

func printStruct(stru interface{}) {
	value := reflect.ValueOf(stru)
	elem := value.Elem()
	for i := 0; i < elem.NumField(); i++ {
		switch elem.Field(i).Kind() {
		case reflect.String, reflect.Int32, reflect.Int64:
			fmt.Printf("%v: %v\n", elem.Type().Field(i).Name, elem.Field(i))
		default:
		}

	}
}
