SHELL := /bin/bash

profile = ddouglas

tidy:
	go mod tidy && go mod vendor

deploy-api:
	@echo ${lambda} for profile ${profile}
	go mod tidy && go mod vendor
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOWORK=off go build -o ./dist/${lambda} ./functions/${lambda}/*.go
	chmod 777 ./dist/${lambda}
	zip -j ./dist/${lambda}.zip ./dist/${lambda}
	aws-vault exec ${AWS_VAULT_DEFAULT_EXEC_PROFILE} -- aws lambda update-function-code --function-name ${lambda}_handler --zip-file "fileb://./dist/${lambda}.zip"
	rm -r ./dist/${lambda}*
