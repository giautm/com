variable "region" {
  type    = string
  default = "us-central1"
}

variable "domain" {
  type    = string
  default = "giautm.dev"
}

# The region in which to put the SQL DB: it is currently configured to use
# PostgreSQL.
# https://cloud.google.com/sql/docs/postgres/locations
variable "db_location" {
  type    = string
  default = "us-central1"
}

# The region for the networking components.
# https://cloud.google.com/compute/docs/regions-zones
variable "network_location" {
  type    = string
  default = "us-central1"
}

# The region for the key management service.
# https://cloud.google.com/kms/docs/locations
variable "kms_location" {
  type    = string
  default = "us-central1"
}

# The location for the app engine; this implicitly defines the region for
# scheduler jobs as specified by the cloudscheduler_location variable but the
# values are sometimes different (as in the default values) so they are kept as
# separate variables.
# https://cloud.google.com/appengine/docs/locations
variable "appengine_location" {
  type    = string
  default = "us-central"
}

# The cloudscheduler_location MUST use the same region as appengine_location but
# it must include the region number even if this is omitted from the
# appengine_location (as in the default values).
variable "cloudscheduler_location" {
  type    = string
  default = "us-central1"
}


# The region in which cloudrun jobs are executed.
# https://cloud.google.com/run/docs/locations
variable "cloudrun_location" {
  type    = string
  default = "us-central1"
}

# The location holding the storage bucket for exported files.
# https://cloud.google.com/storage/docs/locations
variable "storage_location" {
  type    = string
  default = "US"
}

variable "project" {
  type = string
}

variable "cloudsql_tier" {
  type    = string
  default = "db-custom-8-30720"

  description = "Size of the Cloud SQL tier. Set to db-custom-1-3840 or a smaller instance for local dev."
}

variable "cloudsql_disk_size_gb" {
  type    = number
  default = 256

  description = "Size of the Cloud SQL disk, in GB."
}

variable "service_environment" {
  type    = map(map(string))
  default = {}

  description = "Per-service environment overrides."
}

variable "vpc_access_connector_max_throughput" {
  type    = number
  default = 1000

  description = "Maximum provisioned traffic throughput in Mbps"
}

terraform {
  required_providers {
    google      = "~> 3.24"
    google-beta = "~> 3.24"
    null        = "~> 2.1"
    random      = "~> 2.2"
  }
}
