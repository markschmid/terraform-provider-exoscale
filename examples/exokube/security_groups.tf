resource "exoscale_security_group" "exokube" {
  name = "exokube"
  description = "Minikube default Security Group"
}

resource "exoscale_security_group_rule" "exokube_ping" {
  type = "INGRESS"
  description = "Ping"
  security_group_id = "${exoscale_security_group.exokube.id}"
  protocol = "ICMP"
  icmp_type = 8
  icmp_code = 0
  cidr = "0.0.0.0/0"
}

resource "exoscale_security_group_rule" "exokube_ssh" {
  type = "INGRESS"
  description = "SSH"
  security_group_id = "${exoscale_security_group.exokube.id}"
  protocol = "TCP"
  start_port = 22
  end_port = 22
  cidr = "0.0.0.0/0"
}

resource "exoscale_security_group_rule" "exokube_api_server" {
  type = "INGRESS"
  description = "Kubernetes API Server"
  security_group_id = "${exoscale_security_group.exokube.id}"
  protocol = "TCP"
  start_port = 6443
  end_port = 6443
  cidr = "0.0.0.0/0"
}
