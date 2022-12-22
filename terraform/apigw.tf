resource "aws_apigatewayv2_api" "biller" {
  name          = "biller-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_headers = ["*"]
    allow_methods = [
      "GET", "POST", "PATCH",
      "DELETE", "HEAD"
    ]
    expose_headers = ["*"]
    allow_origins  = ["*"]
  }
}

resource "aws_apigatewayv2_domain_name" "biller" {
  domain_name = "api.${local.default_domain}"

  domain_name_configuration {
    certificate_arn = aws_acm_certificate_validation.biller.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}


resource "aws_apigatewayv2_stage" "biller" {
  api_id      = aws_apigatewayv2_api.biller.id
  name        = "$default"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.apigateway_access_logs.arn
    format = jsonencode({
      "requestId"               = "$context.requestId",
      "ip"                      = "$context.identity.sourceIp",
      "requestTime"             = "$context.requestTime",
      "httpMethod"              = "$context.httpMethod",
      "routeKey"                = "$context.routeKey",
      "status"                  = "$context.status",
      "protocol"                = "$context.protocol",
      "responseLength"          = "$context.responseLength",
      "integrationErrorMessage" = "$context.integrationErrorMessage"
    })
  }
}

resource "aws_apigatewayv2_api_mapping" "biller" {
  api_id      = aws_apigatewayv2_api.biller.id
  domain_name = aws_apigatewayv2_domain_name.biller.id
  stage       = aws_apigatewayv2_stage.biller.id
}

resource "aws_apigatewayv2_route" "options_proxy" {
  api_id    = aws_apigatewayv2_api.biller.id
  route_key = "OPTIONS /{proxy+}"
}
