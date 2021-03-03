
GO      ?= go
DOCKER  ?= docker

name := auth-demo
binary := target/$(name)
package := github.com/simia-tech/$(name)/cmd/$(name)

$(binary):
	CGO_ENABLED=0 GOOS=linux $(GO) build -a -o $@ $(package)

image: $(binary)
	$(DOCKER) build -t $(name) .

clean:
	rm -f $(binary)
