#!/bin/env bash


rm -rf dist
mkdir dist 
go mod tidy && go mod vendor
go build -o ./dist ./functions/...

for lambda in $(ls ./dist); do
    chmod 777 ./dist/${lambda}
	zip -j ./dist/${lambda}.zip ./dist/${lambda}
    rm ./dist/$lambda
done
sleep 5
aws-vault exec ddouglas -- go run ./cmd/deploy-functions/*.go
rm -rf dist