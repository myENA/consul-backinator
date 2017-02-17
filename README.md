[![Go Report Card](https://goreportcard.com/badge/github.com/myENA/consul-backinator)](https://goreportcard.com/report/github.com/myENA/consul-backinator) [![Mozilla Public License](https://img.shields.io/badge/license-MPL-blue.svg)](https://www.mozilla.org/MPL/) [![GoDoc](https://godoc.org/github.com/myENA/consul-backinator/common?status.svg)](https://godoc.org/github.com/myENA/consul-backinator/common) [![Docker Automated Build](https://img.shields.io/docker/automated/jrottenberg/ffmpeg.svg)](https://hub.docker.com/r/myena/consul-backinator/)

# consul-backinator

## Summary

Flexible Consul KV pair backup and restore tool with a few unique features
including ACL token backup and restoration.
This was written for and tested in a production environment.

## Note

There was a potentially breaking change in the operation of this tool on
May 09, 2016 when the different operations were broken out into sub commands
to simplify the flag listings.  The functionality remains the same and previous
backup files are __not affected__.  However the command options and structure have
changed and any scripts which embedded this tool will need to be updated.

## Key Features

* Written in Golang using the official Consul API
* No limits on the number of keys that can be backed up or restored
* Backup files are written as gzip compressed and AES256 encrypted JSON data
* Data integrity validation via HMAC-SHA256 signature of the raw data
* Optional path transformation (path replacement) on backup and/or restore
* Clean well documented code that's simple to follow
* Direct AWS/S3 support for backup and restoration of KVs and ACLs

## Installing

OS X/macOS Homebrew users ...

```
brew install consul-backinator
```

Users with a proper Go environment ...

```
go get -u github.com/myENA/consul-backinator
```

Developers that wish to take advantage of vendoring and other options ...

```
git clone https://github.com/myENA/consul-backinator.git
cd consul-backinator
make
```

To build as a Docker container ...

```
git clone https://github.com/myENA/consul-backinator.git
cd consul-backinator
make docker
```

To use the latest container from Docker Hub ...

```
docker pull myena/consul-backinator
docker run myena/consul-backinator
```

See [DOCKER.md](DOCKER.md) for some Docker use cases.

## Security

Releases after 1.4 will be accompanied by a GPG signed SHA256SUM file.

````
ENA R&D Team (ENA R&D code signing key) <r&d@ena.com> 73F17750
````

The public key [62859FAD5BAEA13C3839D5053CA59EE673F17750](http://pgp.key-server.io/search/0x62859FAD5BAEA13C3839D5053CA59EE673F17750)
is available via public servers and [gpg/r&d@ena.com.asc](gpg/r&d@ena.com.asc) in this repository.

When verifying releases please note that only the SHA256SUM file is signed with the GPG key.
The following process should be more than satisfactory to verify the authenticity of any release.

### Verifying Releases

Using the key downloaded from the public server or contained in this repository you
should import the gpg key into your local key store.

```
gpg --import "r&d@ena.com.asc"
```

Next, with the release archive, shasum and signature file downloaded you should
verify integrity of the checksum file using the GPG signature.

```
gpg --verify consul-backinator-1.5-SHA256SUMS.sig consul-backinator-1.5-SHA256SUMS
```

Finally, ensure the downloaded archive matches the verified checksums file.
Only downloaded archives in the current directory will be checked.

```
shasum -a 256 -c consul-backinator-1.5-SHA256SUMS
```

If the above steps completed without error you have a verified release!

## Usage

### Summary

```
ahurt$ ./consul-backinator --help
usage: consul-backinator [--version] [--help] <command> [<args>]

Available commands are:
    backup     Perform a backup operation
    dump       Dump a backup file
    restore    Perform a restore operation

```

### Backup Options

| Option      | Description |
|-------------|-------------|
| `file`      | The backup file target.  The signature will be the same with a `.sig` extension appended.  The default names are `consul.bak` and `consul.bak.sig`
| `key`       | The passphrase used for data encryption and signature generation.  The default string `password` will be used if none specified.  This should be a secure pseudo random string.
| `nokv`      | Do not attempt to backup kv data.  This only makes sense if also passing the `acls` option below.
| `acls`      | Optional backup filename for acl tokens.
| `transform` | Optional argument that affects the key paths written to the backup file.  See the transformation notes below for more information.
| `prefix`    | Optional argument that specifies the starting point for the backup tree.  The default prefix is the root `/` prefix.  To perform a partial tree backup specify a prefix.

### Restore Options

| Option   | Description |
|----------|-------------|
| `file`   | The source file. The default is `consul.bak`
| `key`    | The passphrase used for data decryption and signature validation.  This must match the key used when the backup was created.
| `nokv`   | Do not attempt to restore kv data.  This only makes sense if also passing the `acls` option below.
| `acls`   | Optional source filename for acl tokens.
| `delete` | Optionally delete all keys under the specified prefix prior to restoring the backup file.  The default is false.
| `prefix` | The prefix with the `delete` option.  The default is `/` root.  __THIS WILL DELETE ALL DATA IN YOUR KEYSTORE__ if not changed when using `-delete`.

### Shared Consul Options (backup/restore)

| Option   | Description |
|----------|-------------|
| `addr`            | Optional consul agent address and port.  The default is read from the `CONSUL_HTTP_ADDR` environment variable if specified or set to `127.0.0.1:8500`.
| `scheme`          | Optional scheme `http` or `https` used when connecting to the consul agent.  The default is set to `https` if the `CONSUL_HTTP_SSL` environment variable is set to `true` otherwise the default is `http`.
| `dc`              | Optional datacenter specification.  The default value is the datacenter of the agent to which you are connecting.
| `token`           | Optional consul access token.  The default value is read from the `CONSUL_HTTP_TOKEN` environment variable if specified.
| `ca-cert`         | Optional path to a PEM encoded CA cert file.  This may also be a certificate bundle (concatenation of CA certificates).
| `client-cert`     | Optional path to a PEM encoded client certificate.  This certificate must match the client key.
| `client-key`      | Optional path to an unencrypted PEM encoded private key. This key should obviously match the client cert.  Passing this or `client-cert` alone will probably not work.
| `tls-skip-verify` | Optional bool for verifying a TLS certificate.  This is a very clear security risk and is not reccomended.  This option alone has the same affect as the `CONSUL_HTTP_SSL_VERIFY` environment variable.

### Dump Options

| Option   | Description |
|----------|-------------|
| `file`   | The source file.  The default `consul.bak` will be used if not specified.
| `key`    | The passphrase for the backup file to be dumped.  The default is `password` if not passed.
| `plain`  | Optionally dump and decode a reduced set of information omitting metadata, timestamps and other lesser used information.
| `acls`   | Specified file is an ACL token backup file.  This option is only relevant if using the `plain` option above.

## Transformations

Transformations are simple string operations and will affect the path anywhere
there is a match.  For example, passing `-transform="foo,bar"` would rewrite
`/apple/foo/key` => `/apple/bar/key` as well as `/orange/thing/foo/key` => `/orange/thing/bar/key`.
To avoid potential errors in transformations you should always use the most exact path possible.
Using the previous example if you only wanted to affect keys under `apple` you should pass
`-transform="apple/foo,apple/bar"` to prevent other paths from being modified inadvertently.

## S3 Support

Support for S3 is implemented by passing an S3 URI to the standard ```-file``` option.  The full format for the URI is as follows:

```
s3://access-key:secret-key@my-bucket/path/to/object?region=us-east-1&endpoint=my-s3gw:9000&secure=false
```

The minimal URI when using environment variables would be: ```s3://my-bucket/path/object```

The current URI parsing accepts `s3://` and `s3n://` scheme prefixes.

This table describes all the S3 URI options and corresponding environment variables.

| Paramater    | Environment             | Required    | Description            | Default          |
|--------------|-------------------------|-------------|------------------------|------------------|
| `access-key` | `AWS_ACCESS_KEY_ID`     | yes         | Your S3/AWS access key |                  |
| `secret-key` | `AWS_SECRET_ACCESS_KEY` | yes         | Your S3/AWS secret key |                  |
| `region`     | `AWS_REGION`            | no          | Your S3/AWS region     | us-east-1        |
| `endpoint`   |                         | no          | Optional endpoint      | s3.amazonaws.com |
| `secure`     |                         | no          | Optional secure flag   | true             |

## Example

```
ahurt$ ./consul-backinator backup -key=superSecretStuff
2016/05/09 17:14:11 [Success] Backed up 289 keys from / to consul.bak
Keep your backup (consul.bak) and signature (consul.bak.sig) files in a safe place.
You will need both to restore your data.
```

```
ahurt$ ls -la *.sig *.bak
-rw-------  1 ahurt  staff  11167 May  9 17:14 consul.bak
-rw-------  1 ahurt  staff     44 May  9 17:14 consul.bak.sig
```

## Thanks

* The [HashiCorp](https://github.com/hashicorp) folks for the excellent Consul service discovery daemon and the excellent embedded API and golang package.  In addition to some very nice references back to this utility in a few issues.
* [Adam Avilla](https://github.com/hekaldama) for his time and contributions to the Docker portions of the build.
* All the other contributors, testers and stargazers.
