SHELL := /bin/bash

profile = ddouglas
vault = aws-vault exec ${profile}
distributionID = E1MXR3LWGU7J6O

tidy:
	go mod tidy && go mod vendor

deploy-api:
	@echo ${lambda} for profile ${profile}
	go mod tidy && go mod vendor
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOWORK=off go build -o ./dist/${lambda} ./functions/${lambda}/*.go
	chmod 777 ./dist/${lambda}
	zip -j ./dist/${lambda}.zip ./dist/${lambda}
	${vault} -- aws lambda update-function-code --function-name ${lambda}_handler --zip-file "fileb://./dist/${lambda}.zip"
	rm -r ./dist/${lambda}*

deploy-frontend:
	@echo Deployment Frontend Code
	cd frontend && npm run build
	${vault} -- aws s3 sync --delete frontend/dist/ s3://biller-frontend-us-east-1
	${vault} -- go run ./cmd/create-invalidation/*.go ${distributionID}

plan:
	${vault} -- terraform -chdir=terraform plan

apply:
	${vault} -- terraform -chdir=terraform apply

bastion:
	${vault} -- go run ./cmd/scale-up-asg/main.go