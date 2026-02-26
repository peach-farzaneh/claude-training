# RDS PostgreSQL Module Example
#
# Demonstrates:
# - Multi-AZ for production, single-AZ for staging
# - Parameter groups tuned for Go connection pooling
# - Automated backups
# - Security groups restricted to EKS nodes

variable "db_name" {
  description = "Name of the default database."
  type        = string
  default     = "appdb"
}

variable "db_username" {
  description = "Master username for the database."
  type        = string
  default     = "appuser"
  sensitive   = true
}

variable "db_password" {
  description = "Master password for the database. Use AWS Secrets Manager in production."
  type        = string
  sensitive   = true
}

variable "db_instance_class" {
  description = "RDS instance class."
  type        = string
  default     = "db.t3.medium"
}

variable "eks_security_group_id" {
  description = "Security group ID of the EKS nodes allowed to access RDS."
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs for the DB subnet group."
  type        = list(string)
}

# --- Subnet Group ---

resource "aws_db_subnet_group" "main" {
  name       = "${var.project}-${var.environment}"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.project}-${var.environment}"
  }
}

# --- Security Group ---

resource "aws_security_group" "rds" {
  name_prefix = "${var.project}-rds-"
  vpc_id      = module.vpc.vpc_id
  description = "Allow PostgreSQL access from EKS nodes only."

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [var.eks_security_group_id]
    description     = "PostgreSQL from EKS nodes"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound"
  }

  lifecycle {
    create_before_destroy = true
  }
}

# --- Parameter Group (tuned for Go) ---

resource "aws_db_parameter_group" "main" {
  family = "postgres16"
  name   = "${var.project}-${var.environment}"

  parameter {
    name  = "max_connections"
    value = "200"
  }

  parameter {
    name  = "idle_in_transaction_session_timeout"
    value = "30000"
  }

  parameter {
    name  = "statement_timeout"
    value = "60000"
  }

  parameter {
    name  = "log_min_duration_statement"
    value = "1000"
  }
}

# --- RDS Instance ---

resource "aws_db_instance" "main" {
  identifier = "${var.project}-${var.environment}"
  engine     = "postgres"

  engine_version = "16.4"
  instance_class = var.db_instance_class

  allocated_storage     = 20
  max_allocated_storage = 100
  storage_encrypted     = true

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  multi_az            = var.environment == "production"
  deletion_protection = var.environment == "production"
  skip_final_snapshot = var.environment == "staging"

  backup_retention_period = 7
  backup_window           = "03:00-04:00"
  maintenance_window      = "sun:04:00-sun:05:00"

  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name
  parameter_group_name   = aws_db_parameter_group.main.name

  tags = {
    Name        = "${var.project}-${var.environment}"
    Environment = var.environment
  }
}

# --- Outputs ---

output "rds_endpoint" {
  description = "RDS instance endpoint (host:port)."
  value       = aws_db_instance.main.endpoint
}

output "rds_database_name" {
  description = "Name of the default database."
  value       = aws_db_instance.main.db_name
}

output "rds_connection_string" {
  description = "PostgreSQL connection string (password excluded)."
  value       = "postgres://${var.db_username}@${aws_db_instance.main.endpoint}/${var.db_name}"
  sensitive   = true
}
