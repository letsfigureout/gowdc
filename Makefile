BINDIR=cmd/gowdc
BIN=main.go
OUTPUT=gowdc
DOCKERTAG=${OUTPUT}:latest
REGISTRYTAG=gcr.io/go-apis/gowdc:latest

run:
	go run ${BINDIR}/${BIN}
docker:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ${OUTPUT} ${BINDIR}/${BIN} 
	docker build . -t ${DOCKERTAG}
	docker tag ${DOCKERTAG} ${REGISTRYTAG}
docker-push:
	docker push ${REGISTRYTAG}
clean:
	rm gowdc
