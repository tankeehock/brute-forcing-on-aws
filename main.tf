terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.38.0"
    }
  }
  required_version = ">= 0.15.1"
}

provider "aws" {
  profile = var.profile
  region  = var.region
}

resource "aws_iam_policy" "ec2_role_policy" {
  name = "policy-618033"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = ["sns:Publish", "s3:GetObject"]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_role" "ec2_role" {
  name = "ec2_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = "1"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })
  managed_policy_arns = [aws_iam_policy.ec2_role_policy.arn]
  tags = {
    tag-key = var.tag
  }
}

resource "aws_iam_instance_profile" "ec2_instance_role" {
  name = "ec2_instance_role"
  role = aws_iam_role.ec2_role.name
}

resource "random_id" "s3_bucket_suffix" {
  byte_length = 8
}
resource "aws_s3_bucket" "s3_bucket" {
  bucket = "s3-bucket-${random_id.s3_bucket_suffix.hex}"
  acl    = "private"
  tags = {
    Name = var.tag
  }
}

resource "aws_s3_bucket_object" "app_object" {
  bucket = aws_s3_bucket.s3_bucket.id
  key    = "app"
  source = var.application_file_path
  etag   = filemd5(var.application_file_path)
  tags = {
    Name = var.tag
  }
}

resource "aws_vpc" "main_vpc" {
  cidr_block           = "10.0.0.0/16"
  instance_tenancy     = "default"
  enable_dns_hostnames = true
  tags = {
    Name = var.tag
  }
}

resource "aws_subnet" "public_subnet" {
  vpc_id                  = aws_vpc.main_vpc.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = "true"
  availability_zone       = var.avaliability_zone
  tags = {
    Name = var.tag
  }
}

resource "aws_internet_gateway" "main_igw" {
  vpc_id = aws_vpc.main_vpc.id
  tags = {
    Name = var.tag
  }
}

resource "aws_route_table" "public_subnet_rt" {
  vpc_id = aws_vpc.main_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main_igw.id
  }
  tags = {
    Name = var.tag
  }
}

resource "aws_route_table_association" "public_subnet_rt_association" {
  subnet_id      = aws_subnet.public_subnet.id
  route_table_id = aws_route_table.public_subnet_rt.id
}

resource "aws_security_group" "main_sg_ssh_allowed" {
  vpc_id = aws_vpc.main_vpc.id
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.whitelisted_ip_cidr_ssh]
  }
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags = {
    Name = var.tag
  }
}

resource "aws_vpc_endpoint" "sts_endpoint" {
  vpc_id              = aws_vpc.main_vpc.id
  service_name        = "com.amazonaws.${var.region}.sts"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true
  security_group_ids = [
    aws_security_group.main_sg_ssh_allowed.id,
  ]
}

resource "aws_vpc_endpoint_subnet_association" "sts_endpoint_association" {
  vpc_endpoint_id = aws_vpc_endpoint.sts_endpoint.id
  subnet_id       = aws_subnet.public_subnet.id
}

resource "tls_private_key" "default_key_pair" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "default_generated_key" {
  key_name   = var.key_pair_name
  public_key = tls_private_key.default_key_pair.public_key_openssh
}

resource "aws_instance" "ec2_instances" {
  ami                    = var.ec2_ami
  instance_type          = var.ec2_instance_type
  subnet_id              = aws_subnet.public_subnet.id
  vpc_security_group_ids = [aws_security_group.main_sg_ssh_allowed.id]
  key_name               = aws_key_pair.default_generated_key.id
  iam_instance_profile   = aws_iam_instance_profile.ec2_instance_role.name
  count                  = var.ec2_instances_count
  user_data              = <<EOT
		#! /bin/bash
    aws s3 cp "s3://${aws_s3_bucket_object.app_object.bucket}/${aws_s3_bucket_object.app_object.key}" /home/ec2-user/app
    chmod +x /home/ec2-user/app
    /home/ec2-user/app -region ${var.region} -format "${var.target_access_key}" -n ${var.number_of_characters_for_brute_force} -secret "${var.target_secret_key}" -node-index ${count.index} -number-of-nodes ${var.ec2_instances_count} -random -workers ${var.number_of_workers} -phone-number "${var.phone_number}"
	EOT
  tags = {
    Name = "${var.tag}-${count.index}"
  }
  depends_on = [aws_vpc_endpoint.sts_endpoint, aws_vpc_endpoint_subnet_association.sts_endpoint_association]
}