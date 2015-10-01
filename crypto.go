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
func (c *config) buildKey() []byte {
	sum := sha256.Sum256([]byte(c.cryptKey))
	return sum[:]
}

// read a compressed/encrypted backup file
func (c *config) readFile() ([]byte, error) {
	var in *os.File            // input file
	var gzReader *gzip.Reader  // compressed reader
	var iv [aes.BlockSize]byte // initialization vector
	var cb cipher.Block        // cipher block interface
	var outBytes *bytes.Buffer // output buffer
	var err error              // general error handler

	// open source file
	if in, err = os.Open(c.inFile); err != nil {
		return nil, err
	}

	// close when done
	defer in.Close()

	// wrap reader
	if gzReader, err = gzip.NewReader(in); err != nil {
		return nil, err
	}

	// close when done
	defer gzReader.Close()

	// init cipher block
	if cb, err = aes.NewCipher(c.buildKey()); err != nil {
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
func (c *config) writeFile(data []byte) error {
	var out *os.File           // destination file
	var gzWriter *gzip.Writer  // compressed writer
	var iv [aes.BlockSize]byte // initialization vector
	var cb cipher.Block        // cipher block interface
	var err error              // general error handler

	// open destination file and create/overwite if neeeded
	// and ensure it's only accessible by the current executer
	if out, err = os.OpenFile(c.outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return err
	}

	// close when done
	defer out.Close()

	// wrap writer
	gzWriter = gzip.NewWriter(out)

	// close when done
	defer gzWriter.Close()

	// init cipher block
	if cb, err = aes.NewCipher(c.buildKey()); err != nil {
		return err
	}

	// copy data to destination file encrypting and compressing along the way
	_, err = io.Copy(&cipher.StreamWriter{
		S: cipher.NewOFB(cb, iv[:]), W: gzWriter,
	}, bytes.NewReader(data))

	// return last error
	return err
}
