package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

func debugJson(jsonfile string) error {
	var vmap map[string]interface{}
	data, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &vmap)
	if err != nil {
		return err
	}
	for k, v := range vmap {
		fmt.Fprintf(os.Stdout, "[%s] type [%s]\n", k, reflect.ValueOf(v).Type().String())
	}
	return nil
}

func main() {
	for _, c := range os.Args[1:] {
		debugJson(c)
	}
	return
}
