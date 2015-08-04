all: container

.PHONY: container clean

anomi:
	docker run -it -v $(CURDIR):/gopath/src/github.com/anominet/anomi -w /gopath/src/github.com/anominet/anomi kiasaki/alpine-golang sh -c "go get && go build -a -tags netgo -installsuffix netgo ."

container: anomi
	docker build -t anomi/api .

clean:
	rm anomi || true
