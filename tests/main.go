package main

import (
	"fmt"
	"regexp"
)

func main() {

	fmt.Println("Test")

	patten := "https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z0-9]{2,}(\.[a-zA-Z0-9]{2,})(\.[a-zA-Z0-9]{2,}?"

	r, e := regexp.Compile(patten)

	if e == nil {
		m := r.Match([]byte("http://hello.com"))
		fmt.Println(m)

	} else {
		fmt.Println(e)
	}

}
