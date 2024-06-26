variable "aws_access_key" {
  description = "The IAM public access key"
}

variable "aws_secret_key" {
  description = "IAM secret Key"
}

variable "aws_region" {
  description = "The AWS region"
}

variable "ec2_task_execution_role_name" {
  description = "ECS task execution role name"
  default     = "myEcsTaskExecutionRole"
}

variable "ecs_auto_scale_role_name" {
  description = "ECS auto scale role name"
  default     = "myEcsAutoScaleRole"
}

variable "az_count" {
  description = "Number of AZs to cover in the given region"
  default     = "2"
}

variable "app_image" {
  description = "Docker image to run in ECS"
  default     = "docker.io/akhilbisht798/streamer"
}

variable "app_port" {
  description = "PORT exposed by the docker image"
  default     = 3000
}

variable "alb_port" {
  description = "PORT exposed to the alb"
  default = 80
}

variable "app_count" {
  description = "Number of container to run"
  default     = 3
}

variable "health_check_path" {
  default = "/"
}

variable "fargate_cpu" {
  description = "Fargate instance CPU unit to provision"
  default     = "1024"
}

variable "fargate_memory" {
  description = "Fargate instance memory to provision"
  default     = "2048"
}

variable "cookie_name" {
  description = "cookie name for sticky session"
  default     = "cookieName"
}
