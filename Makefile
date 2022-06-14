IMAGE_NAME = codewars-hook

test:
	go test -shuffle=on -v ./...

docker:
	docker build -t $(IMAGE_NAME) .

