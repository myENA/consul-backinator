package common

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"
)

// WriteBytes writes an encrypted/compressed stream to an io.Writer
func writeBytes(out io.Writer, key string, data []byte) error {
	var gzWriter *gzip.Writer  // compressed writer
	var iv [aes.BlockSize]byte // initialization vector
	var cb cipher.Block        // cipher block interface
	var err error              // general error holder

	// init cipher block
	if cb, err = aes.NewCipher(hashKey(key)); err != nil {
		return err
	}

	// init encrypted writer
	encWriter := &cipher.StreamWriter{
		S: cipher.NewOFB(cb, iv[:]),
		W: out,
	}

	// close when done
	defer encWriter.Close()

	// wrap encrypted writer
	gzWriter = gzip.NewWriter(encWriter)

	// close when done
	defer gzWriter.Close()

	// copy data to destination file compressing and encrypting along the way
	_, err = io.Copy(gzWriter, bytes.NewReader(data))

	// return copy error
	return err
}

// WriteFile writes an encrypted/compressed file and signature
func WriteFile(fname, key string, data []byte) error {
	var out *os.File // destination file
	var err error    // general error holder

	// open destination file and create/overwite if neeeded
	// and ensure it's only accessible by the current executer
	if out, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return err
	}

	// close when done
	defer out.Close()

	// encrypt/compress data and write to output file
	if err = writeBytes(out, key, data); err != nil {
		return err
	}

	// write signature and return
	return writeFileChecksum(fname+".sig", key, data)
}
