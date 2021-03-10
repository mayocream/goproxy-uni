.PHONY: all test clean

test:
	go run main.go -c config.test.yaml

docker:
	sudo docker build -t gu:test .