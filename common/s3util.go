package common

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// Exported error messages
var (
	ErrS3MissingBucketKey = errors.New("missing S3 bucket or object key")
	ErrS3UnknownScheme    = errors.New("unknown URI scheme")
)

// s3Info contains the information needed to connect to an S3
// datastore and create or retrieve objects
type s3Info struct {
	awsConfig *aws.Config
	bucket    string
	key       string
}

// isS3 does a very basic check if the given string *could* be an S3 URI
func isS3(s string) bool {
	return strings.HasPrefix(s, "s3://") || strings.HasPrefix(s, "s3n://")
}

// parseS3URI returns a struct containing all the information needed to connect
// to an S3 endpoing and create or retrieve objects.  The data is collected from
// parsing the passed s3uri and environment variables.
func parseS3URI(s3uri string) (*s3Info, error) {
	var info *s3Info // parsed info
	var u *url.URL   // parsed uri
	var err error    // general error holder

	// parse the s3 path
	if u, err = url.Parse(s3uri); err != nil {
		return nil, err
	}

	// check scheme for giggles
	if u.Scheme != "s3" && u.Scheme != "s3n" {
		return nil, ErrS3UnknownScheme
	}

	// init info
	info = &s3Info{awsConfig: aws.NewConfig()}

	// get access key
	if u.User != nil && u.User.Username() != "" {
		if temps, ok := u.User.Password(); ok {
			info.awsConfig.Credentials = credentials.NewStaticCredentials(
				u.User.Username(),
				temps,
				"",
			)
		}
	}

	// get region
	if temps := u.Query().Get("region"); temps != "" {
		info.awsConfig.Region = aws.String(temps)
	}

	// get bucket
	if info.bucket = u.Host; info.bucket == "" {
		return nil, ErrS3MissingBucketKey
	}

	// get object key
	if info.key = u.Path; u.Path == "" || u.Path == "/" {
		return nil, ErrS3MissingBucketKey
	}

	// check for endpoint override
	if temps := u.Query().Get("endpoint"); temps != "" {
		info.awsConfig.Endpoint = aws.String(temps)
	}

	// check for ssl override
	if temps := u.Query().Get("secure"); temps != "" {
		var secure bool // local bool
		if secure, err = strconv.ParseBool(temps); err != nil {
			return nil, err
		}
		// update config
		info.awsConfig.DisableSSL = aws.Bool(secure)
	}

	// return populated struct
	return info, nil
}
