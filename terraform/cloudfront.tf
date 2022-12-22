resource "aws_cloudfront_origin_access_control" "biller" {
  name                              = "biller-cloudfront-oac"
  description                       = "OAC Policy from Biller Cloudfront"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_distribution" "biller" {
  origin {
    domain_name              = aws_s3_bucket.biller.bucket_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.biller.id
    origin_id                = local.s3_origin_id
  }

  aliases = [local.default_domain]

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "Cloudfront that serves the Biller Frontend"
  default_root_object = "index.html"

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = local.s3_origin_id

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = aws_acm_certificate_validation.biller.certificate_arn
    minimum_protocol_version = "TLSv1.2_2021"
    ssl_support_method       = "sni-only"
  }
}



