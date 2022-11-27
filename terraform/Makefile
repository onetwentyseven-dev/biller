SHELL := /bin/bash

profile = ddouglas
vault = aws-vault exec ${profile}

plan:
	${vault} -- terraform plan

apply:
	${vault} -- terraform apply