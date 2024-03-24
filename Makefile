build:
	go build -ldflags "-s -w" -o gpustatus .

install: build
	mv gpustatus ${HOME}/.local/bin/gpustatus