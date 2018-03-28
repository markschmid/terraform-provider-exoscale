provider "template" {
  version = "~> 1.0"
}

provider "local" {
  version = "~> 1.1"
}

provider "null" {
  version = "~> 1.0"
}

provider "exoscale" {
  version = "~> 0.9.18"
  key = "${var.key}"
  secret = "${var.secret}"
}
