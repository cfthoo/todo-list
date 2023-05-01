.PHONY: build
build:
	docker build -t todo-list .

.PHONY: test
test:
	go test ./... -cover