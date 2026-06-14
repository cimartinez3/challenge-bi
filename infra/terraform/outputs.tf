output "cluster_endpoint" {
  description = "EKS cluster API endpoint — usar para configurar kubectl"
  value       = aws_eks_cluster.main.endpoint
}

output "cluster_name" {
  description = "EKS cluster name"
  value       = aws_eks_cluster.main.name
}

output "ecr_repository_url" {
  description = "ECR repository URL — usar en el pipeline CI/CD para docker push"
  value       = aws_ecr_repository.novobanco.repository_url
}

output "aurora_endpoint" {
  description = "Aurora writer endpoint — nunca exponerlo públicamente"
  value       = aws_rds_cluster.aurora.endpoint
  sensitive   = true
}

output "aurora_reader_endpoint" {
  description = "Aurora reader endpoint — para queries de solo lectura"
  value       = aws_rds_cluster.aurora.reader_endpoint
  sensitive   = true
}

output "secret_arn" {
  description = "ARN del secreto en Secrets Manager con las credenciales de Aurora"
  value       = aws_secretsmanager_secret.aurora_creds.arn
}

output "irsa_role_arn" {
  description = "ARN del IAM Role para anotar el ServiceAccount de Kubernetes"
  value       = aws_iam_role.novobanco_irsa.arn
}
