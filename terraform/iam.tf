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

resource "aws_iam_role_policy_attachment" "lambda_s3" {
  role       = aws_iam_role.biller_lambda_execution_role.name
  policy_arn = aws_iam_policy.lambda_s3.arn
}

resource "aws_iam_policy" "lambda_s3" {
  name   = "LambdaS3Policy"
  policy = data.aws_iam_policy_document.lambda_s3_policy_document.json
}

data "aws_iam_policy_document" "lambda_s3_policy_document" {
  statement {
    sid    = "ListObjectsInBucket"
    effect = "Allow"
    actions = [
      "s3:ListBucket",
    ]
    resources = [
      aws_s3_bucket.receipts.arn
    ]
  }
  statement {
    sid    = "AllObjectActions"
    effect = "Allow"
    actions = [
      "s3:*Object"
    ]
    resources = [
      "${aws_s3_bucket.receipts.arn}/*"
    ]
  }
}
