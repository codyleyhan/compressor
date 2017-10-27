package compressor

import (
	"io"
	"log"

	"github.com/dgryski/go-bitstream"
)

// Helper functions to deal with working with bitstreams

func readBit(r *bitstream.BitReader) (value bool, EOF bool) {
	bit, err := r.ReadBit()

	if err == io.EOF {
		return false, false
	}

	if err != nil {
		log.Fatal("compressor: Problem reading stream", err.Error())
	}

	return bit == bitstream.One, false
}

func readByte(r *bitstream.BitReader) byte {
	b, err := r.ReadByte()

	if err != nil {
		log.Fatal("compressor: Problem reading stream", err.Error())
	}

	return b
}

func writeByte(w *bitstream.BitWriter, b byte) {
	if err := w.WriteByte(b); err != nil {
		log.Fatal("Unable to encode:", err.Error())
	}
}

func writeBit(w *bitstream.BitWriter, val bitstream.Bit) {
	if err := w.WriteBit(val); err != nil {
		log.Fatal("Unable to encode:", err.Error())
	}
}
