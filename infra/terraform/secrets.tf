resource "aws_secretsmanager_secret" "aurora_creds" {
  name                    = "${var.cluster_name}/aurora-credentials"
  description             = "Aurora PostgreSQL connection credentials for novobanco"
  recovery_window_in_days = 7

  tags = {
    Environment = var.environment
  }
}

resource "aws_secretsmanager_secret_version" "aurora_creds" {
  secret_id = aws_secretsmanager_secret.aurora_creds.id

  secret_string = jsonencode({
    host     = aws_rds_cluster.aurora.endpoint
    port     = 5432
    dbname   = var.db_name
    url = "postgres://${var.db_username}@${aws_rds_cluster.aurora.endpoint}:5432/${var.db_name}?sslmode=require"
  })
}
