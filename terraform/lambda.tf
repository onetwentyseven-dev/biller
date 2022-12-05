locals {
  default_env = {
    "SSM_PREFIX" = "/biller"
    "DB_HOST" : data.aws_rds_cluster.ots_cluster.endpoint
    "DB_USER" : "biller"
    "DB_SCHEMA" : "biller"
    "AUTH_TENANT" : "https://onetwentyseven.us.auth0.com/",
    "AUTH_CLIENT_ID" : "u2hgEu2s28xKcKWw0JRgqlgcA6hRLatk",
    "AUTH_AUDIENCE" : "bill-api-development-resource"
  }
}
