data "aws_iam_policy_document" "irsa_assume" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.eks.arn]
    }

    condition {
      test     = "StringEquals"
      variable = "${replace(aws_iam_openid_connect_provider.eks.url, "https://", "")}:sub"
      values   = ["system:serviceaccount:default:novobanco"]
    }

    condition {
      test     = "StringEquals"
      variable = "${replace(aws_iam_openid_connect_provider.eks.url, "https://", "")}:aud"
      values   = ["sts.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "novobanco_irsa" {
  name               = "${var.cluster_name}-irsa"
  assume_role_policy = data.aws_iam_policy_document.irsa_assume.json

  tags = {
    Environment = var.environment
  }
}

resource "aws_iam_role_policy" "secrets_access" {
  name = "${var.cluster_name}-secrets-access"
  role = aws_iam_role.novobanco_irsa.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["secretsmanager:GetSecretValue"]
      Resource = aws_secretsmanager_secret.aurora_creds.arn
    }]
  })
}
