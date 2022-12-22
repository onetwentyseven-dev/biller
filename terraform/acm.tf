resource "aws_acm_certificate" "biller" {
  domain_name       = local.default_domain
  validation_method = "DNS"

  subject_alternative_names = ["*.${local.default_domain}"]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate_validation" "biller" {
  certificate_arn         = aws_acm_certificate.biller.arn
  validation_record_fqdns = [for record in aws_route53_record.biller_cert_validation : record.fqdn]
}
