provider "google" {
  project = var.project
  region  = var.region
}

# For beta-only resources like secrets-manager
provider "google-beta" {
  project = var.project
  region  = var.region
}

# To generate passwords.
provider "random" {}

data "google_project" "project" {
  project_id = var.project
}

resource "google_project_service" "services" {
  project = data.google_project.project.project_id
  for_each = toset([
    "cloudbuild.googleapis.com",
    "cloudkms.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "cloudscheduler.googleapis.com",
    "compute.googleapis.com",
    "containerregistry.googleapis.com",
    "run.googleapis.com",
    "secretmanager.googleapis.com",
    "servicenetworking.googleapis.com",
    "sql-component.googleapis.com",
    "sqladmin.googleapis.com",
    "storage-api.googleapis.com",
    "vpcaccess.googleapis.com",
  ])
  service            = each.value
  disable_on_destroy = false
}

resource "google_compute_global_address" "private_ip_address" {
  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = "projects/${data.google_project.project.project_id}/global/networks/default"

  depends_on = [
    google_project_service.services["compute.googleapis.com"],
  ]
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = "projects/${data.google_project.project.project_id}/global/networks/default"
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]

  depends_on = [
    google_project_service.services["compute.googleapis.com"],
    google_project_service.services["servicenetworking.googleapis.com"],
  ]
}

resource "google_vpc_access_connector" "connector" {
  project        = data.google_project.project.project_id
  name           = "serverless-vpc-connector"
  region         = var.network_location
  network        = "default"
  ip_cidr_range  = "10.8.0.0/28"
  max_throughput = var.vpc_access_connector_max_throughput

  depends_on = [
    google_project_service.services["compute.googleapis.com"],
    google_project_service.services["vpcaccess.googleapis.com"],
  ]
}

# Build creates the container images. It does not deploy or promote them.
resource "null_resource" "build" {
  provisioner "local-exec" {
    environment = {
      PROJECT_ID = data.google_project.project.project_id
      REGION     = var.cloudrun_location
      SERVICES   = "all"
      TAG        = "initial"
    }

    command = "${path.module}/../scripts/build"
  }

  depends_on = [
    google_project_service.services["cloudbuild.googleapis.com"],
  ]
}

# Grant Cloud Build the ability to deploy images. It does not do so in these
# configurations, but it will do future deployments.
resource "google_project_iam_member" "cloudbuild-deploy" {
  project = data.google_project.project.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${data.google_project.project.number}@cloudbuild.gserviceaccount.com"

  depends_on = [
    google_project_service.services["cloudbuild.googleapis.com"],
  ]
}

locals {
  common_cloudrun_env_vars = [
    {
      name  = "PROJECT_ID"
      value = data.google_project.project.project_id
    },
  ]
}

# # Cloud Scheduler requires AppEngine projects!
# resource "google_app_engine_application" "app" {
#   project     = data.google_project.project.project_id
#   location_id = var.appengine_location
# }

output "project_id" {
  value = data.google_project.project.project_id
}

output "project_number" {
  value = data.google_project.project.number
}

output "region" {
  value = var.region
}

output "db_location" {
  value = var.db_location
}

output "network_location" {
  value = var.network_location
}

output "kms_location" {
  value = var.kms_location
}

output "appengine_location" {
  value = var.appengine_location
}

output "cloudscheduler_location" {
  value = var.cloudscheduler_location
}

output "cloudrun_location" {
  value = var.cloudrun_location
}

output "storage_location" {
  value = var.storage_location
}
