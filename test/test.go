package main

import (
	"github.com/ketchupsalt/gohaml"
	"io/ioutil"
	"fmt"
)

func main() {
	c, _ := ioutil.ReadFile("src/github.com/ketchupsalt/gohaml/test/test.haml")
	
	e, err := gohaml.NewEngine(string(c))
	if err != nil {
		fmt.Printf("error: %v", err)
	} else {
		fmt.Println(e.Render(nil))
	}
}
