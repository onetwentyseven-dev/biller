resource "aws_ssm_parameter" "biller_db_password" {
  name  = "/biller/db_pass"
  value = "ChangeMe1234"
  type  = "SecureString"

  lifecycle {
    ignore_changes = [
      value
    ]
  }
}
