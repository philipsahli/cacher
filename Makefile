.PHONY: proto cacher plugin-cacher-filesystem plugin-cacher-redis plugin-cacher-dynamodb
cacher:
		go build -o cacher .

plugin-cacher-filesystem:
		go build -o cacher-filesystem ./plugin-cacher-filesystem

plugin-cacher-redis:
		go build -o cacher-redis ./plugin-cacher-redis

plugin-cacher-dynamodb:
		go build -o cacher-dynamodb ./plugin-cacher-dynamodb

build: proto cacher plugin-cacher-filesystem plugin-cacher-redis plugin-cacher-dynamodb

proto:
		buf generate

all: build run-redis run-dynamodb

run-redis:
		echo
		KV_PLUGIN=./cacher-redis ./cacher put hello redis
		echo
		KV_PLUGIN=./cacher-redis ./cacher get hello
		echo

run-dynamodb:
		KV_PLUGIN=./cacher-dynamodb ./cacher put hello dynamodb
		echo
		KV_PLUGIN=./cacher-dynamodb ./cacher get hello
