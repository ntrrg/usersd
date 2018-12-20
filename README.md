[![Travis build btatus](https://travis-ci.com/ntrrg/usersd.svg?branch=master)](https://travis-ci.com/ntrrg/usersd)
[![codecov](https://codecov.io/gh/ntrrg/usersd/branch/master/graph/badge.svg)](https://codecov.io/gh/ntrrg/usersd)
[![goreport](https://goreportcard.com/badge/github.com/ntrrg/usersd)](https://goreportcard.com/report/github.com/ntrrg/usersd)
[![GoDoc](https://godoc.org/github.com/ntrrg/usersd/pkg/usersd?status.svg)](https://godoc.org/github.com/ntrrg/usersd/pkg/usersd)
[![BCH compliance](https://bettercodehub.com/edge/badge/ntrrg/usersd?branch=master)](https://bettercodehub.com/results/ntrrg/usersd)
[![Docker Build Status](https://img.shields.io/docker/build/ntrrg/usersd.svg)](https://cloud.docker.com/u/ntrrg/repository/docker/ntrrg/usersd)
[![](https://images.microbadger.com/badges/image/ntrrg/usersd.svg)](https://microbadger.com/images/ntrrg/usersd)

**usersd** is an authentication and authorization microservice.

## Building

* Go 1.11.4 or Docker 18.09

### Development:

#### Mage

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

#### Golint

```sh
go get -u -v golang.org/x/lint/golint
```

#### Gometalinter

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

## Contributing

See the [contribution guide](CONTRIBUTING.md) for more information.

## Acknowledgment

Working on this project I use/used:

* [Debian](https://www.debian.org/)

* [XFCE](https://xfce.org/)

* [st](https://st.suckless.org/)

* [Zsh](http://www.zsh.org/)

* [GNU Screen](https://www.gnu.org/software/screen)

* [Git](https://git-scm.com/)

* [EditorConfig](http://editorconfig.org/)

* [Vim](https://www.vim.org/)

* [GNU make](https://www.gnu.org/software/make/)

* [Chrome](https://www.google.com/chrome/browser/desktop/index.html)

* [Gogs](https://gogs.io/)

* [Github](https://github.com)

* [Docker](https://docker.com)

* [Drone](https://drone.io/)

* [Travis CI](https://travis-ci.org)

* [Go Report Card](https://goreportcard.com)

* [Better Code Hub](https://bettercodehub.com)

* [Codecov](https://codecov.io)

* [Mage](https://magefile.org/)

