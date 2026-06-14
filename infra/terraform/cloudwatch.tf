resource "aws_cloudwatch_log_group" "novobanco_api" {
  name              = "/novobanco/api"
  retention_in_days = 30

  tags = {
    Environment = var.environment
  }
}

resource "aws_cloudwatch_log_group" "eks_cluster" {
  name              = "/aws/eks/${var.cluster_name}/cluster"
  retention_in_days = 7

  tags = {
    Environment = var.environment
  }
}
