package main

import (
	"fmt"
	"log"
	"os"

	"github.com/k0kubun/pp"

	"./myjson"
)

func main() {

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 1000)
	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	str := string(data)

	fmt.Println(str)
	parsed, err := myjson.Parse(str)

	pp.Println(parsed)
	fmt.Println(err)
}
