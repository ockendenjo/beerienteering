resource "aws_s3_bucket" "backend" {
  bucket_prefix = "beerienteering-web-"
}

resource "aws_s3_bucket_policy" "allow_cloudfront" {
  bucket = aws_s3_bucket.backend.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          AWS = aws_cloudfront_origin_access_identity.main.iam_arn
        }
        Action   = "s3:GetObject"
        Resource = "${aws_s3_bucket.backend.arn}/*"
      }
    ]
  })
}

output "web_bucket" {
  value = aws_s3_bucket.backend.id
}
