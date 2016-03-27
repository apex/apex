variable "apex_function_hello" {}

resource "aws_api_gateway_rest_api" "Hello" {
  name = "Hello World API"
  description = "Example of Apex's Terraform integration"
}

resource "aws_api_gateway_resource" "Hello" {
  rest_api_id = "${aws_api_gateway_rest_api.Hello.id}"
  parent_id = "${aws_api_gateway_rest_api.Hello.root_resource_id}"
  path_part = "hello"
}

resource "aws_api_gateway_method" "HelloGet" {
  rest_api_id = "${aws_api_gateway_rest_api.Hello.id}"
  resource_id = "${aws_api_gateway_resource.Hello.id}"
  http_method = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "200" {
  rest_api_id = "${aws_api_gateway_rest_api.Hello.id}"
  resource_id = "${aws_api_gateway_resource.Hello.id}"
  http_method = "${aws_api_gateway_method.HelloGet.http_method}"
  status_code = "200"
  response_models = {
    "application/json" = "Empty"
  }
}

resource "aws_api_gateway_integration_response" "Hello" {
  rest_api_id = "${aws_api_gateway_rest_api.Hello.id}"
  resource_id = "${aws_api_gateway_resource.Hello.id}"
  http_method = "${aws_api_gateway_method.HelloGet.http_method}"
  status_code = "${aws_api_gateway_method_response.200.status_code}"
}

resource "aws_api_gateway_integration" "HelloGet" {
  rest_api_id = "${aws_api_gateway_rest_api.Hello.id}"
  resource_id = "${aws_api_gateway_resource.Hello.id}"
  http_method = "${aws_api_gateway_method.HelloGet.http_method}"
  type = "AWS"
  integration_http_method = "POST" # Must be POST for invoking Lambda function
  credentials = "${aws_iam_role.gateway_invoke_lambda.arn}"
  # http://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
  uri = "arn:aws:apigateway:${var.aws_region}:lambda:path/2015-03-31/functions/${var.apex_function_hello}/invocations"
}

resource "aws_api_gateway_deployment" "dev" {
  depends_on = ["aws_api_gateway_integration.HelloGet"]

  rest_api_id = "${aws_api_gateway_rest_api.Hello.id}"
  stage_name = "dev"
}