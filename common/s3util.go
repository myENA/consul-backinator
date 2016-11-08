package common

import (
	"errors"
	"net/url"
	"os"
	"strconv"
)

// Exported error messages
var (
	ErrS3MissingKey = errors.New("Missing S3 access key.  " +
		"They keys should be passed in the URI or set in the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables.  " +
		"Example: s3://access-key:secret-key@my-bucket/path/to/object")
	ErrS3MissingRegion = errors.New("Missing S3 region.  " +
		"The region should be passed in the URI or set in the AWS_REGION environment variable.  " +
		"Example: s3://my-bucket/path/to/object?region=us-east-1")
	ErrS3MissingBucketPath = errors.New("Missing S3 bucket or path.  " +
		"The bucket and path should be passed in the URI specification.  " +
		"Example: s3://my-bucket/path/to/object")
	ErrS3UnknownScheme    = errors.New("Unknown scheme in S3 URI - please use an 's3://' or 's3n://' scheme")
	ErrCreateUnknownError = errors.New("Failed to create bucket on S3 datastore.") // This shouldn't happen
)

// S3Info contains the information needed to connect to an S3
// datastore and create or retrieve objects
type S3Info struct {
	accessKey string
	secretKey string
	region    string
	bucket    string
	path      string
	endpoint  string
	secure    bool
}

// GetS3Info returns a struct containing all the information needed to connect
// to an S3 endpoing and create or retrieve objects.  The data is collected from
// parsing the passed s3uri and environment variables.
func GetS3Info(s3uri string) (*S3Info, error) {
	var info *S3Info // parsed info
	var u *url.URL   // parsed uri
	var err error    // general error holder

	// parse the s3 path
	if u, err = url.Parse(s3uri); err != nil {
		return nil, err
	}

	// very basic validation test
	if u.Scheme != "s3" && u.Scheme != "s3n" {
		return nil, ErrS3UnknownScheme
	}

	// init info
	info = new(S3Info)

	// get access key
	if u.User != nil && u.User.Username() != "" {
		info.accessKey = u.User.Username()
	} else {
		// check environment
		if info.accessKey = os.Getenv("AWS_ACCESS_KEY_ID"); info.accessKey == "" {
			return nil, ErrS3MissingKey
		}
	}

	// get secret key
	if u.User != nil {
		var ok bool // context sensitive validation holder
		if info.secretKey, ok = u.User.Password(); !ok {
			info.secretKey = ""
		}
	}

	// check secret key
	if info.secretKey == "" {
		// check environment
		if info.secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY"); info.secretKey == "" {
			return nil, ErrS3MissingKey
		}
	}

	// get region
	if info.region = u.Query().Get("region"); info.region == "" {
		if info.region = os.Getenv("AWS_REGION"); info.region == "" {
			return nil, ErrS3MissingRegion
		}
	}

	// get bucket
	if info.bucket = u.Host; info.bucket == "" {
		return nil, ErrS3MissingBucketPath
	}

	// get path
	if info.path = u.Path; u.Path == "" || u.Path == "/" {
		return nil, ErrS3MissingBucketPath
	}

	// check for endpoint override
	if info.endpoint = u.Query().Get("endpoint"); info.endpoint == "" {
		info.endpoint = "s3.amazonaws.com"
	}

	// check for ssl override
	if str := u.Query().Get("secure"); str != "" {
		if info.secure, err = strconv.ParseBool(str); err != nil {
			return nil, err
		}
	} else {
		info.secure = true
	}

	// return populated struct
	return info, nil
}
