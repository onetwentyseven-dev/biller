terraform {
  backend "s3" {
    region         = "us-east-1"
    bucket         = "ots-terraform-state-us-east-1"
    key            = "biller.tfstate"
    dynamodb_table = "ots-terraform-state-us-east-1"
  }
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "aws" {
  region = var.region
  default_tags {
    tags = {
      Repository = "https://github.com/onetwentyseven-dev/biller"
    }
  }
}

locals {
  lambda_vpc_access_policy = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

variable "workspace" {
  type = string
}

variable "region" {
  type = string
}

data "aws_vpc" "onetwentyseven" {
  tags = {
    "Name" = "OneTwentySeven"
  }
}


data "aws_subnets" "app_subnets" {
  tags = {
    "app" = true
  }
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

data "aws_rds_cluster" "ots_cluster" {
  cluster_identifier = "ots-serverless-cluster"
}

