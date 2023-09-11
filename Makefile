.PHONY: build

APPNAME := rssas

build:
	docker buildx build --platform linux/amd64 -t maxzhirnov/rssas:latest .
