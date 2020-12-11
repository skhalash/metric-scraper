  
BINARY_NAME=scraper
IMAGE_NAME=skhalashatsap/blackbox-scraper

all: build build-image push-image

build: deps
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) .

build-image: build
	docker build -t $(IMAGE_NAME) .

push-image: build-image
	docker push $(IMAGE_NAME)

clean:
	rm -f $(BINARY_NAME)

run:
	helm upgrade blackbox-scraper ./chart/blackbox-scraper --install

deps:
	go get ./...
