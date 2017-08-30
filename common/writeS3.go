package common

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// write writes an encrypted/compressed object and signature to an S3 datastore
func (info *s3Info) write(key string, data []byte) error {
	var s3Client *s3.S3     // aws s3 client
	var buf *bytes.Buffer   // data buffer
	var err error           // general error holder
	var awsErr awserr.Error // aws framework error
	var ok bool             // assert check

	// init s3 client
	s3Client = s3.New(session.Must(session.NewSession(info.awsConfig)))

	// attempt to create bucket
	if _, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(info.bucket),
	}); err != nil {
		// ignore non-fatal creation errors
		if awsErr, ok = err.(awserr.Error); ok {
			if awsErr.Code() != s3.ErrCodeBucketAlreadyExists &&
				awsErr.Code() != s3.ErrCodeBucketAlreadyOwnedByYou &&
				awsErr.Code() != "AccessDenied" {
				// not something we catch - return the error
				return err
			}
		} else {
			// other failure - return the error
			return err
		}
	}

	// init byte buffer
	buf = new(bytes.Buffer)

	// populate byte buffer with encrypted/compressed data
	if err = writeBytes(buf, key, data); err != nil {
		return err
	}

	// upload data object
	if _, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(info.bucket),
		Key:    aws.String(info.key),
		Body:   bytes.NewReader(buf.Bytes()),
	}); err != nil {
		return err
	}

	// reset and reuse our buf
	buf.Reset()

	// calculate data checksum
	if err = writeChecksum(buf, key, data); err != nil {
		return err
	}

	// upload signature object
	if _, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(info.bucket),
		Key:    aws.String(info.key + ".sig"),
		Body:   bytes.NewReader(buf.Bytes()),
	}); err != nil {
		return err
	}

	// all good
	return nil
}
