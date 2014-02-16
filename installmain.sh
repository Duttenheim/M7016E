export GOPATH="$(pwd)"/code
export GOBIN=$GOPATH/bin
go install $GOPATH/src/$1.go