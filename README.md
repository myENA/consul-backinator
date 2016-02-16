# consul-backinator

## Summary

Flexible Consul KV pair backup and restore tool with a few unique features.
This was written for and tested in a production environment.

## Key Features

* Written in Golang using the official Consul API
* No limits on the number of keys that can be backed up or restored
* Backup files are written as gzip compressed and AES256 encrypted JSON data
* Data integrity validation via HMAC-SHA256 signature of the raw data
* Optional path transformation (path replacement) on backup and/or restore
* Clean well documented code that's simple to follow

## Installing

With a proper Go environment simply run:

```
go get github.com/myENA/consul-backinator
```

## Usage

```
ahurt$ ./consul-backinator -h
Usage of ./consul-backinator:
  -addr string
      Optional consul instance address and port ("127.0.0.1:8500")
  -backup
      Trigger backup operation
  -dc string
      Optional consul datacenter label for backup and restore
  -delete
      Delete all keys under the destination prefix before restore
  -dump
      Dump backup file contents to stdout and exit when used with the -restore option
  -file string
      File for backup and restore operations (default "consul.bak")
  -key string
      Passphrase used for data encryption and signature validation (default "password")
  -prefix string
      Optional prefix from under which all keys will be fetched (default "/")
  -restore
      Trigger restore operation
  -scheme string
      Optional consul instance scheme ("http" or "https")
  -token string
      Optional consul token to access the target cluster
  -transform string
      Optional folder path transformation (oldPath,newPath...)
```

## Example

```
ahurt$ ./consul-backinator -addr=10.16.0.36:8500 -backup -key=SuperSecretStuff
2015/10/02 15:01:59 [Success] Backed up 63 keys from 10.16.0.36:8500/ to consul.bak
Keep your backup (consul.bak) and signature (consul.bak.sig) files in a safe place.
You will need both to restore your data.
```

```
ahurt$ ls -la *.sig *.bak
-rw-------+ 1 ahurt  staff  1901 Oct  2 15:01 consul.bak
-rw-------+ 1 ahurt  staff    44 Oct  2 15:01 consul.bak.sig
```
