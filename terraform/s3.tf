resource "aws_s3_bucket" "receipts" {
  bucket = "biller-receipts-${var.workspace}"
}
