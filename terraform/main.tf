terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
  }
}

variable "trustchain_namespace" {
  type    = string
  default = "trustchain-system"
}

# Example deployment of a Rego policy ConfigMap
resource "kubernetes_config_map" "trustchain_policy" {
  metadata {
    name      = "trustchain-global-policy"
    namespace = var.trustchain_namespace
  }

  data = {
    "policy.rego" = <<-EOT
      package trustchain.admission

      default allow = false

      allow {
          input.verified == true
      }
    EOT
  }
}
