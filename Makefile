all: container

.PHONY: container

container:
	go build -a -tags netgo -installsuffix netgo .
	docker build -t anomi-api .
