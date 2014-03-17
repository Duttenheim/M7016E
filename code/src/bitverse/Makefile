
build:
	go build
	go build ./dht
	cd main; go build -o bitverse
	mv main/bitverse $(GOPATH)bin

test:
	go test .
