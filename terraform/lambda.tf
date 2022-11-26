locals {
  default_env = {
    "SSM_PREFIX" = "/biller"
    "DB_HOST" : data.aws_rds_cluster.ots_cluster.endpoint
    "DB_USER" : "biller"
    "DB_SCHEMA" : "biller"
  }
}
