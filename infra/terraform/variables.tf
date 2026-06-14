variable "aws_region" {
  description = "AWS region where all resources will be created"
  type        = string
  default     = "us-east-1"
}

variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
  default     = "novobanco"
}

variable "db_name" {
  description = "Aurora PostgreSQL database name"
  type        = string
  default     = "novobanco"
}

variable "db_username" {
  description = "Aurora PostgreSQL master username"
  type        = string
  default     = "novobanco"
}

variable "environment" {
  description = "Deployment environment (dev, staging, prod)"
  type        = string
  default     = "prod"
}
