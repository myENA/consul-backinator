package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io"
	"os"
)

// ensure 32 byte hashed key
func buildKey(s string) []byte {
	sum := sha256.Sum256([]byte(s))
	return sum[:]
}

// read a compressed/encrypted backup file
func readFile(src string, key []byte) ([]byte, error) {
	var inFile *os.File        // input file
	var gzReader *gzip.Reader  // compressed reader
	var iv [aes.BlockSize]byte // initialization vector
	var cb cipher.Block        // cipher block interface
	var outBytes *bytes.Buffer // output buffer
	var err error              // general error handler

	// open source file
	if inFile, err = os.Open(src); err != nil {
		return nil, err
	}

	// close when done
	defer inFile.Close()

	// wrap reader
	if gzReader, err = gzip.NewReader(inFile); err != nil {
		return nil, err
	}

	// close when done
	defer gzReader.Close()

	// init cipher block
	if cb, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	// init output
	outBytes = new(bytes.Buffer)

	// copy data to destination file encrypting and compressing along the way
	_, err = io.Copy(outBytes,
		&cipher.StreamReader{
			S: cipher.NewOFB(cb, iv[:]),
			R: gzReader})

	// return last error state
	return outBytes.Bytes(), err
}

// write an encrypted/compressed backup file
func writeFile(dst string, src, key []byte) error {
	var outFile *os.File       // destination file
	var gzWriter *gzip.Writer  // compressed writer
	var iv [aes.BlockSize]byte // initialization vector
	var cb cipher.Block        // cipher block interface
	var err error              // general error handler

	// open destination file and create/overwite if neeeded
	// and ensure it's only accessible by the current executer
	if outFile, err = os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return err
	}

	// close when done
	defer outFile.Close()

	// wrap writer
	gzWriter = gzip.NewWriter(outFile)

	// close when done
	defer gzWriter.Close()

	// init cipher block
	if cb, err = aes.NewCipher(key); err != nil {
		return err
	}

	// copy data to destination file encrypting and compressing along the way
	_, err = io.Copy(&cipher.StreamWriter{
		S: cipher.NewOFB(cb, iv[:]), W: gzWriter,
	}, bytes.NewReader(src))

	// return last error
	return err
}
