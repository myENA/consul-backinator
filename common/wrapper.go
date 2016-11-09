package common

// WriteData writes an encrypted/compressed object and signature
// to a local file or s3 datastore
func WriteData(dest, key string, data []byte) error {
	var info *s3Info // s3 info struct
	var err error    // general error holder

	// basic check
	if isS3(dest) {
		// parse destination as s3 uri and validate
		if info, err = parseS3URI(dest); err != nil {
			return err
		}
		// attempt to write to s3 source
		return info.write(key, data)
	}
	// still going ... attempt file
	return writeFile(dest, key, data)
}

// ReadData reads an encrypted/compressed file or
// S3 datastore object and validates checksums
func ReadData(src, key string) ([]byte, error) {
	var info *s3Info // s3 info struct
	var err error    // general error holder

	// basic check
	if isS3(src) {
		// parse source as s3 uri and validate
		if info, err = parseS3URI(src); err != nil {
			return nil, err
		}
		// attempt to write to s3 source
		return info.read(key)
	}
	// still going ... attempt file
	return readFile(src, key)
}
