repo=$${REPO:-damianjaniszewski/gcp-gcs-perf-test}
version=$${VERSION:-0.1.3}
tag=$(version)

build:
	go build -v gcp-gcs-perf-test.go config.go init.go
	# CGO_ENABLED=0 go build -v -a -tags 'static netgo' -ldflags '-w' gcp-gcs-perf-test.go config.go init.go
	
run:
	go build -v gcp-gcs-perf-test.go config.go init.go
	# CGO_ENABLED=0 go build -v -a -tags 'static netgo' -ldflags '-w' gcp-gcs-perf-test.go config.go init.go
	./gcp-gcs-perf-test

build-container:
	tar -czv -f context.tar.gz ./Dockerfile ./go.mod ./go.sum ./*.go
	docker rmi $(repo):$(tag) $(repo):latest || true
	docker build -t $(repo):$(tag) -t $(repo):latest - < context.tar.gz
	docker push $(repo):$(tag)
	docker push $(repo):latest
	rm context.tar.gz || true

ver:
	@echo $(repo):$(tag)
