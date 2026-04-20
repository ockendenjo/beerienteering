variable "lambda_bucket" {
  type = string
}

variable "permissions_boundary_arn" {
  type = string
}

variable "env" {
  type = string
}

variable "preview_key" {
  type = string
}

variable "live_object_key" {
  type    = string
  default = ""
}

variable "demo_object_key" {
  type    = string
  default = ""
}
