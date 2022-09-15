test:
	go test -v -cover ./...

server:
	go run main.go

serverfromdocker:
	docker build -t gogo_space .
	docker run -p 8080:8080 gogo_space

.PHONY: test server serverfromdocker