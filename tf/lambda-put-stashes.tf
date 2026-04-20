module "lambda_put_stashes" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  s3_bucket                = data.aws_s3_bucket.build_artifacts.id
  s3_object_key            = local.lambda_manifest["put-stashes"]
  name                     = "beerienteering-put-stashes"
  aws_env                  = var.env
  permissions_boundary_arn = var.permissions_boundary_arn
  project_name             = "beerienteering"

  environment = {
    BUCKET_NAME     = aws_s3_bucket.backend.id
    LIVE_OBJECT_KEY = var.live_object_key
  }
}

resource "aws_lambda_permission" "put_stashes_allow_apig" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_put_stashes.arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.my_api.execution_arn}/*"
}

module "iam_s3_put_stashes" {
  source      = "github.com/ockendenjo/tfmods//iam-s3"
  bucket_arn  = aws_s3_bucket.backend.arn
  role_id     = module.lambda_put_stashes.role_id
  allow_write = true
}
