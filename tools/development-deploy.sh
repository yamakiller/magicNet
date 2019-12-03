#!/bin/bash

#file="/etc/profile"
file="/etc/profile_test"
GOLANG_ENV=$1
echo "setup golang development environment"

if [ ! -n "$1" ]; then
    echo "===>Setting environment variables<==="
    sed -i '$a export PATH=$PATH;$GOLANG_ENV/bin' $file
    sed -i '$a export GOPATH=$GOLANG_ENV'
    source /etc/profile
    echo "===>Environment variable setting completed<==="
fi
echo "Start installing development library"
git clone https://github.com/golang/tools  ${GOPATH}/src/golang.org/x/tools
git clone https://github.com/golang/lint   ${GOPATH}/src/golang.org/x/lint
git clone https://github.com/golang/crypto ${GOPATH}/src/golang.org/x/crypto
git clone https://github.com/golang/image  ${GOPATH}/src/golang.org/x/image
git clone https://github.com/golang/net    ${GOPATH}/src/golang.org/x/net
git clone https://github.com/golang/sys    ${GOPATH}/src/golang.org/x/sys
git clone https://github.com/golang/text   ${GOPATH}/src/golang.org/x/text
git clone https://github.com/ramya-rao-a/go-outline ${GOPATH}/src/github.com/ramya-rao-a/go-outline
git clone https://github.com/acroca/go-symbols ${GOPATH}/src/github.com/acroca/go-symbols
git clone https://github.com/josharian/impl ${GOPATH}/src/github.com/josharian/impl
git clone https://github.com/rogpeppe/godef ${GOPATH}/src/github.com/rogpeppe/godef
git clone https://github.com/sqs/goreturns ${GOPATH}/src/github.com/sqs/goreturns
git clone https://github.com/cweill/gotests ${GOPATH}/src/github.com/cweill/gotests
git clone https://github.com/newhook/go-symbols ${GOPATH}/src/github.com/newhook/go-symbols
git clone https://github.com/dropbox/gogoprotobuf ${GOPATH}/src/github.com/dropbox/gogoprotobuf
curpath=$(cd "$(dirname "$0")";pwd)
cd ${GOPATH}/src
go install github.com/ramya-rao-a/go-outline
go install github.com/acroca/go-symbols
go install golang.org/x/tools/cmd/guru
go install golang.org/x/tools/cmd/gorename
go install github.com/josharian/impl
go install github.com/rogpeppe/godef
go install github.com/sqs/goreturns
go install github.com/golang/lint/golint
go install github.com/cweill/gotests/gotests
go install github.com/newhook/go-symbols
go get -u -v github.com/nsf/gocode
go get -u -v github.com/lukehoban/go-find-references
go get -u -v github.com/tpng/gopkgs
go get -u -v github.com/newhook/go-symbols
go install github.com/dropbox/gogoprotobuf/protoc-gen-gogoslick
go get github.com/google/uuid
go get gopkg.in/oauth2.v3
go get github.com/json-iterator/go
go get github.com/nats-io/go-nats-streaming
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/robertkrimen/otto
cd $curpath
git clone https://github.com/sirupsen/logrus %GOPATH%/src/github.com/sirupsen/logrus
git clone https://github.com/x-cray/logrus-prefixed-formatter %GOPATH%/src/github.com/x-cray/logrus-prefixed-formatter
git clone https://github.com/go-sql-driver/mysql %GOPATH%/src/github.com/go-sql-driver/mysql
git clone https://github.com/mongodb/mongo-go-driver %GOPATH%/src/go.mongodb.org/mongo-driver
git clone https://github.com/gomodule/redigo %GOPATH%/src/github.com/gomodule/redigo
git clone https://github.com/golang/freetype %GOPATH%/src/github.com/golang/freetype
git clone https://github.com/golang/protobuf %GOPATH%/src/github.com/golang/protobuf
git clone https://github.com/gorilla/websocket %GOPATH%/src/github.com/gorilla/websocket
git clone https://github.com/grpc/grpc-go %GOPATH%/src/google.golang.org/grpc
git clone https://github.com/google/go-genproto %GOPATH%/src/google.golang.org/genproto
git clone https://github.com/yamakiller/mgolua  %GOPATH%/src/github.com/yamakiller/mgolua

echo "installation is complete"
echo "setup golang development environment complate"