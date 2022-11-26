resource "aws_iam_role" "biller_lambda_execution_role" {
  name               = "BillerLambdaExecutionRole"
  assume_role_policy = data.aws_iam_policy_document.lambda_execution_role_policy_doc.json

}

data "aws_iam_policy_document" "lambda_execution_role_policy_doc" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }
    effect = "Allow"
  }
}

resource "aws_iam_role_policy_attachment" "lambda" {
  role       = aws_iam_role.biller_lambda_execution_role.id
  policy_arn = local.lambda_vpc_access_policy

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_iam_role_policy_attachment" "lambda_ssm" {
  role       = aws_iam_role.biller_lambda_execution_role.name
  policy_arn = aws_iam_policy.lambda_ssm.arn
}

resource "aws_iam_policy" "lambda_ssm" {
  name   = "LambdaSSMPolicy"
  policy = data.aws_iam_policy_document.lambda_ssm_policy_document.json
}

data "aws_iam_policy_document" "lambda_ssm_policy_document" {
  statement {
    effect = "Allow"
    actions = [
      "ssm:GetParameters",
      "ssm:GetParametersByPath"
    ]
    resources = [
      "arn:aws:ssm:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:parameter/biller/*"
    ]
  }
  statement {
    effect = "Allow"
    actions = [
      "ssm:DescribeParameters",
    ]
    resources = [
      "*"
    ]
  }
}

# resource "aws_iam_role_policy_attachment" "lambda_providers_apigateway" {
#   role       = aws_iam_role.biller_lambda_execution_role.name
#   policy_arn = aws_iam_policy.lambda_providers_apigateway.arn
# }

# resource "aws_iam_policy" "lambda_providers_apigateway" {
#   name   = "LambdaAPIGatewayExecution"
#   policy = data.aws_iam_policy_document.lamdba_providers_apigateway.json
# }

# data "aws_iam_policy_document" "lamdba_providers_apigateway" {
#   statement {
#     effect = "Allow"
#     principals {
#       type        = "Service"
#       identifiers = ["apigateway.amazonaws.com"]
#     }
#     actions = [
#       "lambda:InvokeFunction"
#     ]
#     resources = [
#       aws_lambda_function.providers_handler.arn
#     ]
#     condition {
#       test     = "ArnLike"
#       variable = "AWS:SourceArn"
#       values = [
#         "${aws_apigatewayv2_api.biller.execution_arn}/*/GET/providers"
#       ]
#     }
#   }
# }
