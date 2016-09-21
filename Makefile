all:
	@mkdir -p bin/
	@echo "==> Installing dependencies"
	@go get -d -v ./...

format:
	@echo "==> Formating project ..."
	go fmt ./...

build:
	@echo "==> Building ..."
	@go build -o bin/ensure .

dist:
	@echo "==> Creating executables ..."
	@./dist.sh

container:
	docker build -f Dockerfile -t elcuervo/ensure .

release:
	docker push elcuervo/ensure

create:
	#docker rmi -f ensure-builder
	docker build -t ensure-builder -f Dockerfile.build .
	docker run ensure-builder

clean:
	@rm ensure

test:
	@echo "==> Testing ensure ..."
	@go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs -n1 go get -d
	go test ./...

PHONY: all format test
