test: 
	go test ./...

build:
	cd cmd/audio-len/ && go build .

run:
	cd cmd/audio-len/ && go run . -c config.yml	

build-docker:
	cd deploy && $(MAKE) dbuild	

push-docker:
	cd deploy && $(MAKE) dpush

clean:
	rm -f cmd/audio-len/audio-len
	cd deploy && $(MAKE) clean

