-include version
#####################################################################################
## invoke unit tests
test/unit: 
	go test -v -race ./...
.PHONY: test/unit
## code vet and lint
test/lint: 
	go vet ./...
	go install golang.org/x/lint/golint@latest
	golint -set_exit_status ./...
.PHONY: test/lint

build:
	cd cmd/audio-len/ && go build .

run:
	cd cmd/audio-len/ && go run . -c config.yml	
############################################
git/tag:
	git tag "v$(version)"
git/push-tag:
	git push origin --tags
############################################
clean:
	rm -f cmd/audio-len/audio-len
.PHONY: clean
