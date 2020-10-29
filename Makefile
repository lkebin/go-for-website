all: build

build:
	go build

wire:
	wire gen buhaoyong/app
