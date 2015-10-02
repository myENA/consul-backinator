# consul-backinator

## Summary

Flexible Consul KV pair backup and restore tool with a some unique features.
This was written for and tested in a production environment but is still
a work in progress.  More features will be coming but existing
functionality should not be impacted.

## Key Features

* Written in Golang using the official Consul API
* No limits on the number of keys that can be backed up or restored
* Backup files are written as AES256 encrypted/gzipped compressed JSON
* Optional path transformation (path replacement) on backup and/or restore
* Clean well documented code that's simple to follow

## Installing

With a proper Go environment simply run:

```bash
go get github.com/leprechau/consul-backinator
```

## Usage

```bash
$ ./consul-backinator -h
Usage of ./consul-backinator:
  -addr string
        Consul instance address and port ("127.0.0.1:8500")
  -backup
        Trigger backup operation
  -dc string
        Optional consul datacenter label for backup and restore
  -delete
        Optionally delete all keys under the destination prefix before restore
  -in string
        Input file for restore operations (default "consul.bak")
  -key string
        Encryption key used to secure the destination file on backup and read the input file on restore (default "password")
  -out string
        Output file for backup operations (default "consul.bak")
  -prefix string
        Optional prefix from under which all keys will be fetched (default "/")
  -restore
        Trigger restore operation
  -scheme string
        Optional consul instance scheme ("http" or "https")
  -token string
        Optional consul token to access the target cluster
  -transform string
        Optional path transformation to be applied on backup and restore (oldPath,newPath...)
```
