package common

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"
)

// readBytes reads an encrypted/compressed steam from an io.Reader
// and returns a decoded byte slice
func readBytes(in io.Reader, key string) ([]byte, error) {
	var gzReader *gzip.Reader  // compressed reader
	var iv [aes.BlockSize]byte // initialization vector
	var cb cipher.Block        // cipher block interface
	var outBytes *bytes.Buffer // output buffer
	var err error              // general error handler

	// init cipher block
	if cb, err = aes.NewCipher(hashKey(key)); err != nil {
		return nil, err
	}

	// init encrypted reader
	encReader := &cipher.StreamReader{
		S: cipher.NewOFB(cb, iv[:]),
		R: in,
	}

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

	// return bytes and last error state
	return outBytes.Bytes(), err
}

// ReadFile reads an encrypted/compressed file and validates checksums
func ReadFile(fname, key string) ([]byte, error) {
	var in *os.File     // input file
	var outBytes []byte // output buffer
	var err error       // general error handler

	// open source file
	if in, err = os.Open(fname); err != nil {
		return nil, err
	}

	// close when done
	defer in.Close()

	// read and decode file bytes
	if outBytes, err = readBytes(in, key); err != nil {
		return nil, err
	}

	// validate signature file and data
	if err = validateFileChecksum(fname+".sig", key, outBytes); err != nil {
		return nil, err
	}

	// return bytes and last error state
	return outBytes, err
}
