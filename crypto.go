package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

// bad signature
var ErrBadSignature = errors.New("Signature validation failed.  " +
	"Please check your backup file and associated signature.")

// ensure 32 byte hashed key
func (c *config) buildKey() []byte {
	sum := sha256.Sum256([]byte(c.cryptKey))
	return sum[:]
}

// write signature checksum
func (c *config) writeChecksum(data []byte) error {
	var out *os.File // destination file
	var err error    // general error handler

	// open destination file and create/overwite if neeeded
	// and ensure it's only accessible by the current executer
	if out, err = os.OpenFile(c.outFile+".sig", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return err
	}

	// close when done
	defer out.Close()

	// init encoder
	encoder := base64.NewEncoder(base64.StdEncoding, out)

	// build hmac object
	sig := hmac.New(sha256.New, c.buildKey())

	// compute hash
	sig.Write(data)

	// write hash
	_, err = encoder.Write(sig.Sum(nil))

	// close encoder
	encoder.Close()

	// return last error
	return err
}

// validate checksum
func (c *config) validChecksum(data []byte) error {
	var in *os.File       // input file
	var err error         // general error handler
	var buf *bytes.Buffer // signature buffer

	// open source file
	if in, err = os.Open(c.inFile + ".sig"); err != nil {
		return err
	}

	// close when done
	defer in.Close()

	// init decoder
	decoder := base64.NewDecoder(base64.StdEncoding, in)

	// init buffer
	buf = new(bytes.Buffer)

	// read signature
	if _, err = io.Copy(buf, decoder); err != nil {
		return err
	}

	// build hmac object
	sig := hmac.New(sha256.New, c.buildKey())

	// compute hash
	sig.Write(data)

	// validate signature
	if !hmac.Equal(buf.Bytes(), sig.Sum(nil)) {
		return ErrBadSignature
	}

	// no error - all good
	return nil
}

// read an encrypted/compressed backup file
func (c *config) readBackupFile() ([]byte, error) {
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

	// init cipher block
	if cb, err = aes.NewCipher(c.buildKey()); err != nil {
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

	// validate signature
	if err = c.validChecksum(outBytes.Bytes()); err != nil {
		return nil, err
	}

	// return bytes and last error state
	return outBytes.Bytes(), err
}

// write an encrypted/compressed backup file
func (c *config) writeBackupFile(data []byte) error {
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

	// init cipher block
	if cb, err = aes.NewCipher(c.buildKey()); err != nil {
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
	return c.writeChecksum(data)
}
