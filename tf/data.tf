data "aws_s3_bucket" "build_artifacts" {
  bucket = var.lambda_bucket
}

data "aws_s3_object" "lambda_manifest" {
  bucket = data.aws_s3_bucket.build_artifacts.id
  key    = "lambda_manifests/default.json"
}

locals {
  lambda_manifest = jsondecode(data.aws_s3_object.lambda_manifest.body)
}
