resource "aws_route53_zone" "biller" {
  name = local.default_domain
}

resource "aws_route53_record" "biller" {
  zone_id = aws_route53_zone.biller.zone_id
  name    = local.default_domain
  type    = "A"

  alias {
    zone_id                = aws_cloudfront_distribution.biller.hosted_zone_id
    name                   = aws_cloudfront_distribution.biller.domain_name
    evaluate_target_health = false
  }

}

resource "aws_route53_record" "biller_cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.biller.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = aws_route53_zone.biller.zone_id
}

resource "aws_route53_record" "biller_api" {
  name    = aws_apigatewayv2_domain_name.biller.domain_name
  type    = "A"
  zone_id = aws_route53_zone.biller.zone_id

  alias {
    name                   = aws_apigatewayv2_domain_name.biller.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.biller.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

output "zone_ns_records" {
  value = aws_route53_zone.biller.name_servers
}
