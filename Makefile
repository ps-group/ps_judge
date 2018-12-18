.PHONY: all
all: build

.PHONY: build
build: build_backend

.PHONY: build_backend
build_backend:
	(cd src/backend_service && GOPATH=$(PWD) go build -o backend_service .)
	cp -f bin/backend_service.json src/backend_service/backend_service.json
	docker build -t backend_service src/backend_service
