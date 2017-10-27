package main

import (
	"os"

	"github.com/codyleyhan/compressor"

	"io/ioutil"
	"log"
)

func main() {
	data, err := ioutil.ReadFile("test.txt")

	if err != nil {
		log.Fatal("compressor: Unable to read file:", err.Error())
	}

	f, err := os.Create("test.comp")

	if err != nil {
		log.Fatal("compressor: unable to create file", err.Error())
	}

	compressor.Encode(&data, f)

	if err := f.Close(); err != nil {
		log.Fatal("cannot close", err.Error())
	}

	comp, err := os.Open("test.comp")
	if err != nil {
		log.Fatal("compressor: unable to open file", err.Error())
	}

	defer comp.Close()

	decomp, err := os.Create("test-decom.txt")
	if err != nil {
		log.Fatal("compressor: unable to open file", err.Error())
	}

	defer decomp.Close()

	compressor.Decode(comp, decomp)
}
