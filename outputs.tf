output "ssh_private_key" {
  description = "SSH private key for the EC2 instances"
  value       = tls_private_key.default_key_pair.private_key_pem
  sensitive   = true
}
output "external_ip_addresses" {
  description = "External IP address of the EC2 instances"
  value       = aws_instance.ec2_instances[*].public_ip
}