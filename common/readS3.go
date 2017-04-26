package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// read reads an encrypted/compressed object from an S3 datastore and validates checksums
func (info *s3Info) read(key string) ([]byte, error) {
	var s3Client *s3.S3                           // aws s3 client
	var dataObject, sigObject *s3.GetObjectOutput // fetched objects
	var outBytes []byte                           // output buffer
	var err error                                 // general error holder

	// init s3 client
	s3Client = s3.New(session.Must(session.NewSession(info.awsConfig)))

	// fetch data object adn check error
	if dataObject, err = s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(info.bucket),
		Key:    aws.String(info.key),
	}); err != nil {
		return nil, err
	}

	// read and decode data object
	if outBytes, err = readBytes(dataObject.Body, key); err != nil {
		return nil, err
	}

	// fetch signature object and check error
	if sigObject, err = s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(info.bucket),
		Key:    aws.String(info.key + ".sig"),
	}); err != nil {
		return nil, err
	}

	// validate signature
	if err = validateChecksum(sigObject.Body, key, outBytes); err != nil {
		return nil, err
	}

	// return bytes and last error state
	return outBytes, err
}
