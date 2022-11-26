resource "aws_cloudwatch_log_group" "apigateway_access_logs" {
  name              = "/aws/apigatewayv2/logs"
  retention_in_days = 1
}

resource "aws_cloudwatch_log_group" "providers_handler" {
  name              = "/aws/lambda/${aws_lambda_function.providers_handler.function_name}"
  retention_in_days = 3
}

