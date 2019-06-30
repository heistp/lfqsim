package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func fail(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	os.Exit(1)
}

func printJson(v interface{}) {
	json, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		fail(err)
	}

	fmt.Println(string(json))
}

func main() {
	var c Config
	d := json.NewDecoder(bufio.NewReader(os.Stdin))
	if err := d.Decode(&c); err != nil {
		fail(err)
	}
	printJson(c)

	s := NewSimulator(&c)
	r := s.Run()
	printJson(r)
}
