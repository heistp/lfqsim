package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func printJson(v interface{}) {
	json, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(json))
}

func main() {
	log.SetFlags(0)

	var c Config
	d := json.NewDecoder(bufio.NewReader(os.Stdin))
	if err := d.Decode(&c); err != nil {
		log.Fatalln(err)
	}
	printJson(c)

	r := NewSimulator(&c).Run()

	printJson(r)
}
