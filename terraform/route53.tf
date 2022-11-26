resource "aws_route53_zone" "biller" {
  name = "biller.onetwentyseven.dev"
}


output "zone_ns_records" {
  value = aws_route53_zone.biller.name_servers
}
