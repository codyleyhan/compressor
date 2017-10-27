# Compressor

This package allows for the data compression of ASCII Art.

## Implementation

The data compression is utilizing the lossless Huffman coding algorithm to generate prefixes/codes that are shorter than the standard 8 bits that would be required for each character in standard ASCII.  These codes are based on the the frequency that the character appears in the art, with the most frequent character having the shortest code.  Using this algorithm we can achieve a very high lossless compression.  For the test data, the algorithm was able to compress from 5495 bytes to 1327 bytes which is a compression ratio of ~4. 

In order to serialize the data, I first add some metadata as a header to the beginning of the data.  This metadata includes the number of character to code mappings as 1 byte, the number of total characters that are encoded as 2 bytes, and finally the serialized implementation of the prefix tree that is used in huffman encoding as a variable number of bytes.  I do not encode the frequencies as there is no added information that will aid in decompression by including them. Then the actual data is encoded after the meta data.

![Data Serialization Format](https://docs.google.com/drawings/d/e/2PACX-1vTtJMiJGQZD_GmLkG3fMcQlHTGWLHSs1GY-7qvA6UMhNTcaPitTDO50iGgKQ_p4Vz-par3o3BA67Ka4/pub?w=480&h=360)

## Limitations

Due to the way that I implemented my serialization of metadata, the theoretical limit on the number of characters that can be encoded and decoded is 2^16, but this is far more characters than is asked of as 100 x 100 is only 10k characters.  

The other limitation is that the algorithm does not preform well when the encoded data is relativly small such as one word or a very short sentence.  This is due to the overhead that is needed in serializing the metadata.

## Future

If I had more time I would have done much better unit testing as well as cleaning up some of the code organization.  I would also have liked to make the code much more robust.  Additionally I would like to add the ability to tell when the algorithm will produce poor results and not compress the data if this occurs.

## How to Run

Ensure that the following dependencies are installed:
- Go 1.9
- dep 

Ensure that this repo is in the correct directory of your gopath.

I intentionally left the package very open ended in order for it to be used in multiple use cases.

To run the simple use case ensure you are in the cmd folder:

```bash
# /github.com/codyleyhan/compressor/
dep ensure

cd cmd

go run compressor.go
```

Alternativly the package can be imported and used in another program:

```go
import "github.com/codyleyhan/compressor"

// the following can now be used

func main() {
    compressor.Encode(data, file)

    compressor.Decode(compressedFile, buffer)
}
```