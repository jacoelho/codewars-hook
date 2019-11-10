IMAGE_NAME = codewars-hook

test:
	go test -v ./...

docker:
	docker build -t $(IMAGE_NAME) .

