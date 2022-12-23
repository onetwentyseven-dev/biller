resource "aws_lambda_function" "sheets_handler" {

  filename         = "${path.module}/assets/example/handler.zip"
  source_code_hash = filebase64sha256("${path.module}/assets/example/handler.zip")

  function_name = "sheets_handler"
  role          = aws_iam_role.biller_lambda_execution_role.arn
  handler       = "sheets"
  runtime       = "go1.x"
  timeout       = 30

  vpc_config {
    security_group_ids = data.aws_db_instance.ots_db_instance.vpc_security_groups
    subnet_ids         = data.aws_subnets.lambda_subnets.ids
  }


  environment {
    variables = local.default_env
  }

  lifecycle {
    ignore_changes = [
      source_code_hash
    ]
  }
}


locals {
  sheets_routes = toset([
    "GET /sheets",
    "POST /sheets",
    "GET /sheets/{sheetID}",
    "PATCH /sheets/{sheetID}",
    "DELETE /sheets/{sheetID}",
    "GET /sheets/{sheetID}/entries",
    "POST /sheets/{sheetID}/entries",
    "PATCH /sheets/{sheetID}/entries/{entryID}",
    "DELETE /sheets/{sheetID}/entries/{entryID}"
  ])
}

resource "aws_lambda_permission" "sheets_handler" {
  for_each      = local.sheets_routes
  statement_id  = sha1(each.value)
  function_name = aws_lambda_function.sheets_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.biller.execution_arn}/*/${join("", split(" ", each.value))}"
  action        = "lambda:InvokeFunction"
}

resource "aws_apigatewayv2_route" "sheets_handler" {
  for_each = local.sheets_routes
  depends_on = [
    aws_apigatewayv2_integration.sheets_handler
  ]

  api_id    = aws_apigatewayv2_api.biller.id
  route_key = each.value

  target = "integrations/${aws_apigatewayv2_integration.sheets_handler.id}"

}

resource "aws_apigatewayv2_integration" "sheets_handler" {
  api_id      = aws_apigatewayv2_api.biller.id
  description = "API GatewayV2 Integration for route ${aws_lambda_function.sheets_handler.function_name}"

  integration_type       = "AWS_PROXY"
  integration_method     = "POST"
  integration_uri        = aws_lambda_function.sheets_handler.invoke_arn
  payload_format_version = "2.0"

  lifecycle {
    create_before_destroy = true
  }

}
