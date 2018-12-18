# Build requirements

* Go 1.11.3

For Docker building:

* Docker 18.09

## Development:

### Mage

```sh
$ cd /tmp
```

```sh
$ wget -c 'https://github.com/magefile/mage/releases/download/v1.8.0/mage_1.8.0_Linux-64bit.tar.gz'
```

```sh
$ tar -xf mage_1.8.0_Linux-64bit.tar.gz
```

```sh
$ cp -a /tmp/mage $(go env GOPATH)/bin/
```

```sh
$ cd -
```

### Golint

```sh
go get -u -v golang.org/x/lint/golint
```

### Gometalinter

```sh
$ cd /tmp
```

```sh
$ wget -c 'https://github.com/alecthomas/gometalinter/releases/download/v2.0.11/gometalinter-2.0.11-linux-amd64.tar.gz'
```

```sh
$ tar -xf gometalinter-2.0.11-linux-amd64.tar.gz
```

```sh
$ cp -a $(find gometalinter-2.0.11-linux-amd64/ -type f) $(go env GOPATH)/bin/
```

```sh
$ cd -
```

