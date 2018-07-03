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
	var info *s3Info                // parsed info
	var u *url.URL                  // parsed url
	var accessKey, secretKey string // key holders
	var err error                   // general error holder

	// The `net/url` package does not handle '/' in password.
	// Therefore, we strip out and parse the user/password portion manually.
	// See: https://github.com/myENA/consul-backinator/issues/30
	if strings.Contains(s3uri, "@") {
		var keyStart = strings.Index(s3uri, "://") + 3 // get start
		var keyEnd = strings.Index(s3uri, "@")         // get end
		var keyString = s3uri[keyStart:keyEnd]         // pull out key string
		// check key string
		if strings.Contains(keyString, ":") {
			// split the keys
			var keySplit = strings.Split(keyString, ":")
			// check split
			if len(keySplit) == 2 {
				// set keys
				accessKey = keySplit[0]
				secretKey = keySplit[1]
				// rewrite uri - remove credentials
				s3uri = s3uri[:keyStart] + s3uri[keyEnd+1:]
			}
		}
	}

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

	// check access/secret key
	if accessKey != "" && secretKey != "" {
		info.awsConfig.Credentials = credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		)
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
	if info.key = u.Path; info.key == "" || info.key == "/" {
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
		// update config to disable SSL if secure is false
		// see #43 - this was a terrible name for this paramater
		if secure == false {
			info.awsConfig.DisableSSL = aws.Bool(true)
		}
	}

	//check for pathstyle override
	if temps := u.Query().Get("pathstyle"); temps != "" {
		var pathstyle bool // local bool
		if pathstyle, err = strconv.ParseBool(temps); err != nil {
			return nil, err
		}
		// update config
		info.awsConfig.S3ForcePathStyle = aws.Bool(pathstyle)
	}

	// return populated struct
	return info, nil
}
