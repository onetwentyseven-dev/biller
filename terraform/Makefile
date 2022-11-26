SHELL := /bin/bash

profile = ddouglas
vault = aws-vault exec ${profile}

plan:
	aws-vault exec ${profile} -- terraform plan

apply:
	${vault} -- terraform apply