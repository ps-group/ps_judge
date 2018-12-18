.PHONY: all
all: build

.PHONY: build
build: build_backend build_builder

.PHONY: build_backend
build_backend:
	(cd src/backend_service && GOPATH=$(PWD) go build -o build/backend_service .)
	cp -f bin/backend_service.json src/backend_service/build/backend_service.json
	docker build -t backend_service src/backend_service/build

.PHONY: build_builder
build_builder:
	(cd src/builder_service && GOPATH=$(PWD) go build -o build/builder_service .)
	cp -f bin/builder_service.json src/builder_service/build/builder_service.json
	docker build -t builder_service src/builder_service/build
