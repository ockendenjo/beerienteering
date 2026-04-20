# Define the REST API
resource "aws_api_gateway_rest_api" "my_api" {
  name        = "beerienteering"
  description = "Beerienteering (${var.env})"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "api" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id

  triggers = {
    //Note that a deployment can take a few minutes to propagate even after the `terraform apply` has completed
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.root.path_part,
      aws_api_gateway_method.get_stashes.id,
      aws_api_gateway_method.put_stashes.id,
      module.lambda_get_stashes.invoke_arn,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    aws_api_gateway_integration.get_stashes,
    aws_api_gateway_integration.put_stashes,
  ]
}



resource "aws_api_gateway_stage" "api" {
  deployment_id = aws_api_gateway_deployment.api.id
  rest_api_id   = aws_api_gateway_rest_api.my_api.id
  stage_name    = "api"
}

resource "aws_api_gateway_resource" "root" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  parent_id   = aws_api_gateway_rest_api.my_api.root_resource_id
  path_part   = "stashes"
}

resource "aws_api_gateway_method" "get_stashes" {
  rest_api_id   = aws_api_gateway_rest_api.my_api.id
  resource_id   = aws_api_gateway_resource.root.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "get_stashes" {
  rest_api_id             = aws_api_gateway_rest_api.my_api.id
  resource_id             = aws_api_gateway_resource.root.id
  http_method             = aws_api_gateway_method.get_stashes.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = module.lambda_get_stashes.invoke_arn
}

resource "aws_api_gateway_method" "put_stashes" {
  rest_api_id   = aws_api_gateway_rest_api.my_api.id
  resource_id   = aws_api_gateway_resource.root.id
  http_method   = "PUT"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "put_stashes" {
  rest_api_id             = aws_api_gateway_rest_api.my_api.id
  resource_id             = aws_api_gateway_resource.root.id
  http_method             = aws_api_gateway_method.put_stashes.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = module.lambda_put_stashes.invoke_arn
}
