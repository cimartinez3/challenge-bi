resource "aws_db_subnet_group" "main" {
  name       = "${var.cluster_name}-aurora-subnet-group"
  subnet_ids = aws_subnet.private[*].id

  tags = {
    Name        = "${var.cluster_name}-aurora-subnet-group"
    Environment = var.environment
  }
}

resource "aws_security_group" "aurora" {
  name        = "${var.cluster_name}-aurora-sg"
  description = "Allow PostgreSQL access from EKS nodes only"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.eks_cluster.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.cluster_name}-aurora-sg"
  }
}

resource "aws_rds_cluster" "aurora" {
  cluster_identifier = "${var.cluster_name}-aurora"
  engine             = "aurora-postgresql"
  engine_version     = "15.4"
  database_name      = var.db_name
  master_username    = var.db_username

  manage_master_user_password = true

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.aurora.id]

  skip_final_snapshot = false
  final_snapshot_identifier = "${var.cluster_name}-aurora-final-snapshot"

  backup_retention_period = 7
  preferred_backup_window = "03:00-04:00"

  tags = {
    Environment = var.environment
  }
}

resource "aws_rds_cluster_instance" "writer" {
  identifier         = "${var.cluster_name}-aurora-writer"
  cluster_identifier = aws_rds_cluster.aurora.id
  instance_class     = "db.t3.medium"
  engine             = aws_rds_cluster.aurora.engine
  engine_version     = aws_rds_cluster.aurora.engine_version
  availability_zone  = data.aws_availability_zones.available.names[0]
}

resource "aws_rds_cluster_instance" "reader" {
  identifier         = "${var.cluster_name}-aurora-reader"
  cluster_identifier = aws_rds_cluster.aurora.id
  instance_class     = "db.t3.medium"
  engine             = aws_rds_cluster.aurora.engine
  engine_version     = aws_rds_cluster.aurora.engine_version
  availability_zone  = data.aws_availability_zones.available.names[1]
}
