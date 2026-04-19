module "lambda_get_stashes" {
  source                   = "github.com/ockendenjo/tfmods//lambda"
  s3_bucket                = data.aws_s3_bucket.build_artifacts.id
  s3_object_key            = local.lambda_manifest["get-stashes"]
  name                     = "hbt-beerienteering-get-stashes"
  aws_env                  = var.env
  permissions_boundary_arn = var.permissions_boundary_arn
  project_name             = "beerienteering"

  environment = {
    GO_LIVE_TIME = "2025-07-17T18:58:00+01:00"
    BUCKET_NAME  = aws_s3_bucket.backend.id
    PREVIEW_KEY  = var.preview_key
  }
}

resource "aws_lambda_permission" "get_stashes_allow_apig" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_get_stashes.arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.my_api.execution_arn}/*"
}

resource "aws_iam_role_policy" "get_stashes_s3" {
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid = "AllowS3"
        Action = [
          "s3:List*",
          "s3:GetObject*",
        ]
        Effect = "Allow"
        Resource = [
          aws_s3_bucket.backend.arn,
          "${aws_s3_bucket.backend.arn}/*",
        ]
      },
    ]
  })
  role = module.lambda_get_stashes.role_id
}
