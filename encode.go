package compressor

import (
	"container/heap"
	"encoding/binary"
	"io"

	"github.com/dgryski/go-bitstream"
)

type prefix struct {
	numBits int8
	code    string
}

type encoder struct {
	rawData     *[]byte
	root        *node
	w           *bitstream.BitWriter
	prefixes    map[byte]*prefix
	frequencies map[byte]int
}

// Implements huffman encoding algorithm
func (e *encoder) encode() {
	e.buildFrequencyTable()

	h := buildMinHeap(&e.frequencies)
	e.root = buildPrefixTree(h)

	e.buildPrefixTable(e.root, "")

	writeHeader(e)
	writeData(e.rawData, &e.prefixes, e.w)
}

func (e *encoder) buildFrequencyTable() {
	for _, char := range *e.rawData {
		// returns val or 0
		val := e.frequencies[char]

		e.frequencies[char] = val + 1
	}
}

// builds a min heap based on the frequencies
func buildMinHeap(frequencies *map[byte]int) *minHeap {
	h := make(minHeap, len(*frequencies))

	idx := 0
	for char, frequency := range *frequencies {
		n := node{char, frequency, false, nil, nil}

		h[idx] = &n

		idx++
	}

	heap.Init(&h)

	return &h
}

// generates a unique prefix for each leaf node in the tree
func (e *encoder) buildPrefixTable(n *node, code string) {
	if n == nil {
		return
	}

	// If we are at leaf node add prefix to the table
	if !n.internal {
		p := prefix{int8(len(code)), code}
		e.prefixes[n.char] = &p
	} else {
		e.buildPrefixTable(n.left, code+"0")
		e.buildPrefixTable(n.right, code+"1")
	}
}

func buildPrefixTree(h *minHeap) *node {
	// while we have more than 1 node in the heap
	for {
		if len(*h) <= 1 {
			break
		}

		// grab the two minimum frequency nodes
		left := heap.Pop(h).(*node)
		right := heap.Pop(h).(*node)

		// generate new intenal node
		newNode := node{0, left.frequency + right.frequency, true, left, right}

		heap.Push(h, &newNode)
	}

	return heap.Pop(h).(*node)
}

// Writes header data in the format:
// First byte: number of leaf nodes in tree
// Next 2 bytes: number of characters that will be encoded
// Next n bytes: the serialized tree
func writeHeader(e *encoder) {
	numCharsInTree := len(e.prefixes)

	// number of characters that will be encoded
	numCharsInData := uint16(len(*e.rawData))

	// convert the number to bytes
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, numCharsInData)

	writeByte(e.w, byte(numCharsInTree))
	writeByte(e.w, bs[0])
	writeByte(e.w, bs[1])

	writeTree(e.root, e.w)
}

func writePrefix(p *prefix, w *bitstream.BitWriter) {
	var i int8
	for i = 0; i < p.numBits; i++ {
		if p.code[i] == '0' {
			writeBit(w, bitstream.Zero)
		} else {
			writeBit(w, bitstream.One)
		}
	}
}

// loops over each character and writes the mapped prefix to the writer
func writeData(data *[]byte, prefixes *map[byte]*prefix, w *bitstream.BitWriter) {
	for _, char := range *data {
		p := (*prefixes)[char]
		writePrefix(p, w)
	}
}

//Serialize the tree to the writer
func writeTree(n *node, w *bitstream.BitWriter) {
	if n == nil {
		return
	}

	if n.internal {
		// a zero bit denotes that this node is internal
		writeBit(w, bitstream.Zero)
		writeTree(n.left, w)
		writeTree(n.right, w)
	} else {
		// denote that a leaf node has been reached with a value of one
		writeBit(w, bitstream.One)
		// write the character byte
		writeByte(w, n.char)
	}
}

//Encode compresses the data provided and saves the compressed data to the writer
func Encode(data *[]byte, w io.Writer) {
	e := encoder{}
	e.rawData = data
	e.w = bitstream.NewWriter(w)
	e.frequencies = make(map[byte]int)
	e.prefixes = make(map[byte]*prefix)

	e.encode()
	e.w.Flush(bitstream.One)
}
