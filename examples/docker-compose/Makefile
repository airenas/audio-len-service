-include .env
#####################################################################################
start:
	docker compose up -d

stop:
	docker compose down

test:
	curl -X POST http://localhost:$(PORT)/duration -H 'content-type: multipart/form-data' -F file=@1.mp3

