install: build
	go install

deps: 
	go get ./...
	
build: test
	go build

test:
	go test -race -v -cover

benchmark:
	 go test -bench=. -v -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof
