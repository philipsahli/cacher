# Cacher demo

## Build core

    go build -o cacher .

## Build & run with plugin-cacher-redis 

    go build -o cacher-redis ./plugin-cacher-redis

    docker run --name redis3  -p 6379:6379 -d redis:5.0.3-alpine

    KV_PLUGIN=./cacher-redis ./cacher put host postfinance 
	KV_PLUGIN=./cacher-redis ./cacher get host

## Build & run with plugin-cacher-dynamodb

    go build -o cacher-dynamodb ./plugin-cacher-dynamodb

    docker run -d -p 8000:8000 amazon/dynamodb-local
    aws dynamodb create-table --table-name=MyTable --attribute-definitions AttributeName=key,AttributeType=S --key-schema AttributeName=key,KeyType=HASH --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url http://localhost:8000

    aws dynamodb scan --table-name=MyTable --endpoint-url http://localhost:8000

    KV_PLUGIN=./cacher-dynamodb ./cacher put weather rainy
	KV_PLUGIN=./cacher-dynamodb ./cacher get weather