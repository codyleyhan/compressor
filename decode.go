package compressor

import (
	"encoding/binary"
	"io"
	"log"

	"github.com/dgryski/go-bitstream"
)

type decoder struct {
	root            *node
	w               io.Writer
	r               *bitstream.BitReader
	numCharsEncoded uint16
	numCharsWritten uint16
	numChars        uint8
	numCharsDecoded uint8
}

// decode the meta data and tree in the header
func (d *decoder) decodeHeader() {
	// first byte is the number of leaf nodes
	d.numChars = uint8(readByte(d.r))

	// read in the total number of characters in the encoded data
	buf := make([]byte, 2)
	buf[0] = readByte(d.r)
	buf[1] = readByte(d.r)

	d.numCharsEncoded = binary.LittleEndian.Uint16(buf)

	// deserialize the tree
	d.root = d.createTree()
}

func (d *decoder) decodeData() {
	p := d.root

	// loop over all the bits in the data
	for {
		// read a bit
		one, done := readBit(d.r)
		if done {
			break
		}

		// traverse right for 1
		if one {
			if p.right != nil {
				p = p.right

				// mapping found from code to character
				if !p.internal {
					d.w.Write([]byte{p.char})
					p = d.root

					d.numCharsWritten++
				}
			} else {
				log.Fatal("Problem decoding data")
			}
		} else { // traverse left
			if p.left != nil {
				p = p.left
				if !p.internal {
					d.w.Write([]byte{p.char})
					p = d.root
					d.numCharsWritten++
				}
			} else {
				log.Fatal("Problem decoding data")
			}
		}

		// stop once all characters have been decoded to avoid padding problems
		if d.numCharsEncoded == d.numCharsWritten {
			break
		}
	}
}

// read the bit and create the node, zero indicates an internal node, one indicates a leaf character node
func (d *decoder) createTree() *node {
	if val, _ := readBit(d.r); val {
		return &node{readByte(d.r), -1, false, nil, nil}
	} else if d.numChars != d.numCharsDecoded {
		left := d.createTree()
		right := d.createTree()
		return &node{0, -1, true, left, right}
	}

	return nil
}

func (d *decoder) decode() {
	d.decodeHeader()
	d.decodeData()
}

//Decode takes a compressed reader, decodes the data and outputs the decompressed data to the writer
func Decode(r io.Reader, w io.Writer) {
	d := decoder{}

	d.w = w
	d.r = bitstream.NewReader(r)

	d.decode()
}
