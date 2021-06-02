variable "region" {
  description = "Region for the AWS deployment"
  type        = string
}
variable "avaliability_zone" {
  description = "Avalability zone for the AWS deployment"
  type        = string
}
variable "profile" {
  description = "AWS Configuration Profile"
  type        = string
}
variable "tag" {
  description = "Common tag for easy identification / grouping"
  type        = string
}
variable "ec2_instances_count" {
  description = "Number of EC2 instance(s) to spin up"
  type        = number
}
variable "key_pair_name" {
  description = "Name of the keypair used"
  type        = string
  default     = "default_key_pair"
}
variable "ec2_ami" {
  description = "EC2 AMI Image - Defaults to Amazon Linux 2 AMI (HVM)"
  type        = string
  default     = "ami-0d5eff06f840b45e9"
}
variable "ec2_instance_type" {
  description = "EC2 Instance Type - Defaults to t2.micro"
  type        = string
  default     = "t2.micro"
}
variable "whitelisted_ip_cidr_ssh" {
  description = "Whitelisted IP CIDR for security group(s) (SSH) related to the EC2 instance(s)"
  type        = string
}
variable "phone_number" {
  description = "Phone number to recieve the SMS notification (eg. +65XXXXXXXX)"
  type        = string
}
variable "application_file_path" {
  description = "Relative file path to the application executable"
  type        = string
}
variable "target_access_key" {
  description = "Access key which will be used for the brute-force attack. The access key is formated to the requirements of the binary that will perform the bruteforce (eg. AKIAXXXXX%sXX)"
  type        = string
}
variable "target_secret_key" {
  description = "Secret key which will be used for the brute-force attack"
  type        = string
}
variable "number_of_characters_for_brute_force" {
  description = "Number of characters in the access key to be attacked"
  type        = number
}
variable "number_of_workers" {
  description = "Number of workers spawned from the executable"
  type        = number
}