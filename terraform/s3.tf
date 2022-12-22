resource "aws_s3_bucket" "biller" {
  bucket = "biller-frontend-${var.region}"
}

resource "aws_s3_bucket_acl" "biller" {
  bucket = aws_s3_bucket.biller.bucket
  acl    = "private"
}

resource "aws_s3_bucket_policy" "biller" {
  bucket = aws_s3_bucket.biller.bucket
  policy = data.aws_iam_policy_document.cloudfront_biller_s3_policy.json
}

data "aws_iam_policy_document" "cloudfront_biller_s3_policy" {
  statement {
    sid    = "AllowCloudFrontServicePrincipal"
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["cloudfront.amazonaws.com"]
    }
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.biller.arn}/*"]
    condition {
      test     = "StringEquals"
      variable = "aws:sourceArn"

      values = [
        aws_cloudfront_distribution.biller.arn
      ]
    }
  }
}

# {
#   "Version": "2008-10-17",
#   "Id": "PolicyForCloudFrontPrivateContent",
#   "Statement": [
#       {
#           "Sid": "AllowCloudFrontServicePrincipal",
#           "Effect": "Allow",
#           "Principal": {
#               "Service": "cloudfront.amazonaws.com"
#           },
#           "Action": "s3:GetObject",
#           "Resource": "arn:aws:s3:::biller-frontend-us-east-1/*",
#           "Condition": {
#               "StringEquals": {
#                 "AWS:SourceArn": "arn:aws:cloudfront::816459443219:distribution/E1MXR3LWGU7J6O"
#               }
#           }
#       }
#   ]
# }


resource "aws_s3_bucket" "receipts" {
  bucket = "biller-receipts-${var.workspace}"
}

output "receipts_bucket_name" {
  value = aws_s3_bucket.receipts.bucket
}
