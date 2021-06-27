terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=2.64.0"
    }
  }
}

provider "azurerm" {
  features {}
}

module "akscluster" {
  source       = "github.com/implodingduck/tfmodules//aks"
  cluster_name = var.cluster_name
  location     = var.location
  node_count   = var.node_count
  tags         = {
    owner = "implodingduck"
  }
}
