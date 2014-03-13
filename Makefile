V=0

build:
	export GOPATH=${CURDIR}/code; cd code/src; make
	
clean:
	rm -r -f code/pkg/*
	rm -f code/bin/*
