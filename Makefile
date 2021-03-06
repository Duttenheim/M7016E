V=0

build:
	mkdir -p code/bin
	export GOPATH=${CURDIR}/code; cd code/src; make

setup:
	export GOPATH=${CURDIR}/code; go get github.com/nu7hatch/gouuid; go get code.google.com/p/go.net/websocket; go get github.com/fsouza/go-dockerclient; go get github.com/beevik/ntp
	
clean:
	rm -r -f code/pkg/*
	rm -f code/bin/*
