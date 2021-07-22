data "snowflake_system_get_privatelink_config" "snowflake_private_link" {}

resource "aws_security_group" "snowflake_private_link" {
  vpc_id = var.vpc_id

  ingress {
    from_port   = 80
    to_port     = 80
    cidr_blocks = var.vpc_cidr
    protocol    = "tcp"
  }

  ingress {
    from_port   = 443
    to_port     = 443
    cidr_blocks = var.vpc_cidr
    protocol    = "tcp"
  }
}

resource "aws_vpc_endpoint" "snowflake_private_link" {
  vpc_id              = var.vpc_id
  service_name        = data.snowflake_system_get_privatelink_config.snowflake_private_link.aws_vpce_id
  vpc_endpoint_type   = "Interface"
  security_group_ids  = [aws_security_group.snowflake_private_link.id]
  subnet_ids          = var.subnet_ids
  private_dns_enabled = false
}

resource "aws_route53_zone" "snowflake_private_link" {
  name = "privatelink.snowflakecomputing.com"

  vpc {
    vpc_id = var.vpc_id
  }
}

resource "aws_route53_record" "snowflake_private_link_url" {
  zone_id = aws_route53_zone.snowflake_private_link.zone_id
  name    = data.snowflake_system_get_privatelink_config.snowflake_private_link.account_url
  type    = "CNAME"
  ttl     = "300"
  records = [aws_vpc_endpoint.snowflake_private_link.dns_entry[0]["dns_name"]]
}

resource "aws_route53_record" "snowflake_private_link_oscp_url" {
  zone_id = aws_route53_zone.snowflake_private_link_url.zone_id
  name    = data.snowflake_system_get_privatelink_config.snowflake_private_link.oscp_url
  type    = "CNAME"
  ttl     = "300"
  records = [aws_vpc_endpoint.snowflake_private_link.dns_entry[0]["dns_name"]]
}
