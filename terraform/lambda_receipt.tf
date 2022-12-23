resource "aws_lambda_function" "receipts_handler" {

  filename         = "${path.module}/assets/example/handler.zip"
  source_code_hash = filebase64sha256("${path.module}/assets/example/handler.zip")

  function_name = "receipts_handler"
  role          = aws_iam_role.biller_lambda_execution_role.arn

  handler = "receipts"
  runtime = "go1.x"
  timeout = 30

  vpc_config {
    security_group_ids = data.aws_db_instance.ots_db_instance.vpc_security_groups
    subnet_ids         = data.aws_subnets.lambda_subnets.ids
  }

  environment {
    variables = merge(local.default_env, {
      "RECEIPT_BUCKET" : aws_s3_bucket.receipts.bucket,
    })
  }

  lifecycle {
    ignore_changes = [
      source_code_hash
    ]
  }
}


locals {
  receipts_routes = toset([
    "GET /receipts",
    "POST /receipts",
    "GET /receipts/{receiptID}",
    "PATCH /receipts/{receiptID}",
    "DELETE /receipts/{receiptID}",
    "GET /receipts/{receiptID}/file",
    "POST /receipts/{receiptID}/file",
    "DELETE /receipts/{receiptID}/file",
  ])
}

resource "aws_lambda_permission" "receipts_handler" {
  for_each      = local.receipts_routes
  statement_id  = sha1(each.value)
  function_name = aws_lambda_function.receipts_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.biller.execution_arn}/*/${join("", split(" ", each.value))}"
  action        = "lambda:InvokeFunction"
}

resource "aws_apigatewayv2_route" "receipts_handler" {
  for_each = local.receipts_routes
  depends_on = [
    aws_apigatewayv2_integration.receipts_handler
  ]

  api_id    = aws_apigatewayv2_api.biller.id
  route_key = each.value

  target = "integrations/${aws_apigatewayv2_integration.receipts_handler.id}"

}

resource "aws_apigatewayv2_integration" "receipts_handler" {
  api_id      = aws_apigatewayv2_api.biller.id
  description = "API GatewayV2 Integration for route ${aws_lambda_function.receipts_handler.function_name}"

  integration_type       = "AWS_PROXY"
  integration_method     = "POST"
  integration_uri        = aws_lambda_function.receipts_handler.invoke_arn
  payload_format_version = "2.0"

  lifecycle {
    create_before_destroy = true
  }

}

