# consul-backinator

## Summary

Flexible Consul KV pair backup and restore tool with a few unique features.
This was written for and tested in a production environment.

## Note
There was a breaking change in the operation of this tool on May 09, 2016 when
the different operations were broken out into sub commands to simplify the flag
listing.

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

Or, if you have glide (https://github.com/Masterminds/glide) ...

```
git clone https://github.com/myENA/consul-backinator.git
cd consul-backinator
./build.sh -i
```

## Usage

```
ahurt$ ./consul-backinator --help
usage: consul-backinator [--version] [--help] <command> [<args>]

Available commands are:
    backup     Perform a backup operation
    dump       Dump a backup file
    restore    Perform a backup operation

```

```
ahurt$ ./consul-backinator backup --help
Usage: ./consul-backinator backup [options]

  Performs a backup operation against a consul cluster KV store.

Options:

  -file         Sets the backup filename (default: "consul.bak")
  -key          Passphrase for data encryption and signature validation (default: "password")
  -transform    Optional path transformation (oldPath,newPath...)
  -prefix       Optional prefix from under which all keys will be fetched
  -addr         Optional consul address and port (default: "127.0.0.1:8500")
  -scheme       Optional consul scheme ("http" or "https")
  -dc           Optional consul datacenter
  -token        Optional consul access token

```

```
ahurt$ ./consul-backinator restore --help
Usage: ./consul-backinator restore [options]

  Performs a restore operation against a consul cluster KV store.

Options:

  -file         Source filename (default: "consul.bak")
  -key          Passphrase for data encryption and signature validation (default: "password")
  -transform    Optional path transformation (oldPath,newPath...)
  -delete       Delete all keys under specified prefix prior to restoration (default: false)
  -prefix       Prefix for delete operation
  -addr         Optional consul address and port (default: "127.0.0.1:8500")
  -scheme       Optional consul scheme ("http" or "https")
  -dc           Optional consul datacenter
  -token        Optional consul access token

```

```
ahurt$ ./consul-backinator dump --help
Usage: ./consul-backinator dump [options]

  Dump and optionally decode the contents of a backup file to stdout.

Options:

  -file         Source filename (default: "consul.bak")
  -key          Passphrase for data encryption and signature validation (default: "password")
  -plain        Dump only the key and decoded value

```

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
