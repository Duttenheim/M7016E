V=0

build:
	export GOPATH=${CURDIR}/code; cd code/src; make

setup:
	export GOPATH=${CURDIR}/code; go get github.com/nu7hatch/gouuid; go get code.google.com/p/go.net/websocket
	
clean:
	rm -r -f code/pkg/*
	rm -f code/bin/*
