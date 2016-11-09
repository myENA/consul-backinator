package common

import (
	"bytes"
	"github.com/minio/minio-go"
)

// write writes an encrypted/compressed object and signature to an S3 datastore
func (info *s3Info) write(key string, data []byte) error {
	var mc *minio.Client  // minio s3 client
	var buf *bytes.Buffer // data buffer
	var err error         // general error holder

	// init minio client
	if mc, err = minio.New(info.endpoint, info.accessKey, info.secretKey, info.secure); err != nil {
		return err
	}

	// attempt to create bucket
	if err = mc.MakeBucket(info.bucket, info.region); err != nil {
		var exists bool // context sensitive validator
		// it might already exist
		if exists, err = mc.BucketExists(info.bucket); err != nil {
			return err
		}
		// not exists but no error - don't think this should ever happen
		if !exists {
			return ErrCreateUnknownError
		}
	}

	// init byte buffer
	buf = new(bytes.Buffer)

	// populate byte buffer with encrypted/compressed data
	if err = writeBytes(buf, key, data); err != nil {
		return err
	}

	// upload data object
	if _, err = mc.PutObject(info.bucket, info.path, buf, "application/octet-stream"); err != nil {
		return err
	}

	// reset and reuse our buf
	buf.Reset()

	// calculate data checksum
	if err = writeChecksum(buf, key, data); err != nil {
		return err
	}

	// upload signature object
	if _, err = mc.PutObject(info.bucket, info.path+".sig", buf, "application/octet-stream"); err != nil {
		return err
	}

	// all good
	return nil
}
