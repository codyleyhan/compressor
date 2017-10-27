package compressor

import (
	"bytes"
	"os"
	"testing"
)

var before *bytes.Buffer
var compressed bytes.Buffer
var after bytes.Buffer

// With more time, there would be more thorough unit testing and more tests
func TestMain(m *testing.M) {
	testData := ""

	for i := 0; i < 5000; i++ {
		if i%7 == 0 {
			testData += "a"
		} else if i%3 == 0 {
			testData += "b"
		} else {
			testData += "c"
		}
	}

	before = bytes.NewBufferString(testData)

	returnCode := m.Run()

	os.Exit(returnCode)
}

func TestBuildFrequencyTable(t *testing.T) {
	temp := before.Bytes()

	e := encoder{}
	e.rawData = &temp
	e.frequencies = make(map[byte]int)

	e.buildFrequencyTable()

	if len(e.frequencies) != 3 {
		t.Fatal("The encoder should build a table that has the correct number of characters")
	}

	if e.frequencies[byte(97)] != 715 {
		t.Fatal("The encoder did not count the correct number of occurences for \"a\"")
	}
}

func TestCompression(t *testing.T) {
	temp := before.Bytes()

	Encode(&temp, &compressed)

	if compressed.Len() >= 5000 {
		t.Fatal("The data was not compressed")
	}
}

func TestDecompression(t *testing.T) {
	compressedReader := bytes.NewReader(compressed.Bytes())
	Decode(compressedReader, &after)

	if after.Len() != 5000 {
		t.Fatal("The data was not uncompressed")
	}
}
