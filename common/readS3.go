package common

import (
	"github.com/minio/minio-go"
)

// ReadS3 writes an encrypted/compressed object and signature to an s3 datastore
func (info *S3Info) Read(key string) ([]byte, error) {
	var mc *minio.Client         // minio s3 client
	var dataObject *minio.Object // fetched data object
	var sigObject *minio.Object  // fetched signature object
	var outBytes []byte          // output buffer
	var err error                // general error holder

	// init minio client
	if mc, err = minio.New(info.endpoint, info.accessKey, info.secretKey, info.secure); err != nil {
		return nil, err
	}

	// fetch data object
	if dataObject, err = mc.GetObject(info.bucket, info.path); err != nil {
		return nil, err
	}

	// close when done
	defer dataObject.Close()

	// read and decode data object
	if outBytes, err = readBytes(dataObject, key); err != nil {
		return nil, err
	}

	// fetch signature object
	if sigObject, err = mc.GetObject(info.bucket, info.path+".sig"); err != nil {
		return nil, err
	}

	// validate signature
	if err = validateChecksum(sigObject, key, outBytes); err != nil {
		return nil, err
	}

	// return bytes and last error state
	return outBytes, err
}
