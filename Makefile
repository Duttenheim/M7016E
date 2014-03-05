build:
	cd code/src/bitverse; make
	cd code/src;
	go build docker
	cd code/src/entrypoints; \
		go build -o ${GOPATH}/bin/dockerserver dockerservermain.go; \
		go build -o ${GOPATH}/bin/dockerclient dockerclientmain.go; \
		go build -o ${GOPATH}/bin/bitverseserver bitverseservertest.go; \
		go build -o ${GOPATH}/bin/bitverseclient bitverseclienttest.go; \
	
clean:
	rm -r code/pkg/*
