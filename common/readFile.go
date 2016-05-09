package common

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"
)

// ReadFile reads an encrypted/compressed file and validates checksums
func ReadFile(fname, key string) ([]byte, error) {
	var in *os.File            // input file
	var gzReader *gzip.Reader  // compressed reader
	var iv [aes.BlockSize]byte // initialization vector
	var cb cipher.Block        // cipher block interface
	var outBytes *bytes.Buffer // output buffer
	var err error              // general error handler

	// open source file
	if in, err = os.Open(fname); err != nil {
		return nil, err
	}

	// close when done
	defer in.Close()

	// init cipher block
	if cb, err = aes.NewCipher(hashKey(key)); err != nil {
		return nil, err
	}

	// init encrypted reader
	encReader := &cipher.StreamReader{
		S: cipher.NewOFB(cb, iv[:]),
		R: in}

	// wrap encrypted reader
	if gzReader, err = gzip.NewReader(encReader); err != nil {
		return nil, err
	}

	// close when done
	defer gzReader.Close()

	// init output
	outBytes = new(bytes.Buffer)

	// read data into output buffer decompressing and decrypting along the way
	_, err = io.Copy(outBytes, gzReader)

	// validate signature file and data
	if err = validateChecksum(fname+".sig", key, outBytes.Bytes()); err != nil {
		return nil, err
	}

	// return bytes and last error state
	return outBytes.Bytes(), err
}
