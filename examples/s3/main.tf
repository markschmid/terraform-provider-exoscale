provider "aws" {
  version = "~> 1.12.0"
  access_key = "${var.key}"
  secret_key = "${var.secret}"

  region = "${var.zone}"
  endpoints {
    s3 = "https://sos-ch-dk-2.exo.io"
  }

  # Deactivate the AWS specific behaviours
  #
  # https://www.terraform.io/docs/backends/types/s3.html#skip_credentials_validation
  skip_credentials_validation = true
  skip_get_ec2_platforms = true
  skip_requesting_account_id = true
  skip_metadata_api_check = true
  skip_region_validation = true
}

resource "aws_s3_bucket" "testbucket" {
  bucket = "atestforterraform"
}

resource "aws_s3_bucket_object" "testobj" {
  bucket = "${aws_s3_bucket.testbucket.bucket}"
  key = "some-text.txt"
}
