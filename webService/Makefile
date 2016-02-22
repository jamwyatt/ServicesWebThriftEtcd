
TARGET=webFrontEnd
VERSION=latest
DOCKERTAG=jamwyatt/$(shell echo $(TARGET) | tr A-Z a-z):$(VERSION)

$(TARGET): $(TARGET).go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo $^

docker: $(TARGET) Dockerfile
	docker build --rm=true -t $(DOCKERTAG) .
	touch docker

dockerRun: docker
	docker run --rm=true $(DOCKERTAG)

all: $(TARGET) docker

clean:
	rm -f $(TARGET) docker
	-docker rmi -f $(DOCKERTAG)

