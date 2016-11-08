package common

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"hash"
	"io"
	"os"
)

// ErrBadSignature indicates failed signature validation
var ErrBadSignature = errors.New("Signature validation failed.  " +
	"Please check your backup file and associated signature.")

// hashKey builds a 32 byte hashed key
func hashKey(key string) []byte {
	sum := sha256.Sum256([]byte(key))
	return sum[:]
}

// writeChecksum writes a signature to the given io.Writer
func writeChecksum(out io.Writer, key string, data []byte) error {
	var encoder io.WriteCloser // encoding writer
	var sig hash.Hash          // hash object
	var err error              // general error handler

	// init encoder
	encoder = base64.NewEncoder(base64.StdEncoding, out)

	// close encoder when done
	defer encoder.Close()

	// build hmac object
	sig = hmac.New(sha256.New, hashKey(key))

	// compute hash
	sig.Write(data)

	// write hash
	_, err = encoder.Write(sig.Sum(nil))

	// return write error
	return err
}

// writeFileChecksum writes a signature to a file
func writeFileChecksum(fname, key string, data []byte) error {
	var out *os.File // destination file
	var err error    // general error handler

	// open destination file and create/overwite if neeeded
	// and ensure it's only accessible by the current executer
	if out, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return err
	}

	// close when done
	defer out.Close()

	// return write error
	return writeChecksum(out, key, data)
}

// validateChecksum
func validateChecksum(in io.Reader, key string, data []byte) error {
	var decoder io.Reader // encoding writer
	var sig hash.Hash     // hash object
	var buf *bytes.Buffer // signature buffer
	var err error         // general error handler

	// init decoder
	decoder = base64.NewDecoder(base64.StdEncoding, in)

	// init buffer
	buf = new(bytes.Buffer)

	// read signature
	if _, err = io.Copy(buf, decoder); err != nil {
		return err
	}

	// build hmac object
	sig = hmac.New(sha256.New, hashKey(key))

	// compute hash
	sig.Write(data)

	// validate signature
	if !hmac.Equal(buf.Bytes(), sig.Sum(nil)) {
		return ErrBadSignature
	}

	// no error - all good
	return nil
}

// validateFileChecksum validates a signature file and data
func validateFileChecksum(fname, key string, data []byte) error {
	var in *os.File // input file
	var err error   // general error handler

	// open source file
	if in, err = os.Open(fname); err != nil {
		return err
	}

	// close when done
	defer in.Close()

	// return validation
	return validateChecksum(in, key, data)
}
