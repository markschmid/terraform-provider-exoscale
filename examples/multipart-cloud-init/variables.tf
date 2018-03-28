variable "key" {}
variable "secret" {}
variable "key_pair" {}

variable "hostnames" {
  type = "list"
  default = ["alpha", "beta"]
}

variable "zone" {
  default = "ch-dk-2"
}

variable "template" {
  default = "Linux Ubuntu 17.10 64-bit"
}

