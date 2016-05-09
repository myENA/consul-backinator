package common

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"
)

// WriteFile writes an encrypted/compressed file and signature
func WriteFile(fname, key string, data []byte) error {
	var out *os.File           // destination file
	var gzWriter *gzip.Writer  // compressed writer
	var iv [aes.BlockSize]byte // initialization vector
	var cb cipher.Block        // cipher block interface
	var err error              // general error handler

	// open destination file and create/overwite if neeeded
	// and ensure it's only accessible by the current executer
	if out, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return err
	}

	// close when done
	defer out.Close()

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
	if _, err = io.Copy(gzWriter, bytes.NewReader(data)); err != nil {
		// return copy error
		return err
	}

	// write signature and return
	return writeChecksum(fname+".sig", key, data)
}
