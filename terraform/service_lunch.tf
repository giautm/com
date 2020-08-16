#
# Create and deploy the service
#

resource "google_service_account" "lunch-bot" {
  project      = data.google_project.project.project_id
  account_id   = "viecco-lunch-bot-sa"
  display_name = "Viec.Co Lunch Service"
}

resource "google_service_account_iam_member" "cloudbuild-deploy-lunch-bot" {
  service_account_id = google_service_account.lunch-bot.id
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${data.google_project.project.number}@cloudbuild.gserviceaccount.com"

  depends_on = [
    google_project_service.services["cloudbuild.googleapis.com"],
  ]
}

# resource "google_secret_manager_secret_iam_member" "lunch-bot" {
#   provider = google-beta

#   for_each = toset([
#     "sslcert",
#     "sslkey",
#     "sslrootcert",
#     "password",
#   ])

#   secret_id = google_secret_manager_secret.db-secret[each.key].id
#   role      = "roles/secretmanager.secretAccessor"
#   member    = "serviceAccount:${google_service_account.lunch-bot.email}"
# }

resource "google_cloud_run_service" "lunch-bot" {
  name     = "lunch-bot"
  location = var.cloudrun_location

  template {
    spec {
      service_account_name = google_service_account.lunch-bot.email

      containers {
        image = "gcr.io/${data.google_project.project.project_id}/giautm.dev/viecco/cmd/lunch-bot:initial"

        resources {
          limits = {
            cpu    = "2"
            memory = "1G"
          }
        }

        env {
          name  = "WEBHOOK_BASE_URL"
          value = "https://viecco-lunch.${var.domain}"
        }

        dynamic "env" {
          for_each = local.common_cloudrun_env_vars
          content {
            name  = env.value["name"]
            value = env.value["value"]
          }
        }

        dynamic "env" {
          for_each = lookup(var.service_environment, "lunch-bot", {})
          content {
            name  = env.key
            value = env.value
          }
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" : "3",
        "run.googleapis.com/vpc-access-connector" : "${google_vpc_access_connector.connector.id}"
      }
    }
  }

  depends_on = [
    google_project_service.services["run.googleapis.com"],
    # google_secret_manager_secret_iam_member.lunch-bot,
    null_resource.build,
  ]

  lifecycle {
    ignore_changes = [
      template,
    ]
  }
}

# resource "google_cloud_run_domain_mapping" "lunch-bot" {
#   location = var.cloudrun_location
#   name     = "viecco-lunch.${var.domain}"

#   metadata {
#     namespace = data.google_project.project.project_id
#   }

#   spec {
#     route_name = google_cloud_run_service.lunch-bot.name
#   }

#   depends_on = [
#     google_project_service.services["run.googleapis.com"],
#     google_cloud_run_service.lunch-bot,
#   ]
# }

data "google_iam_policy" "lunch-bot-noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location    = google_cloud_run_service.lunch-bot.location
  project     = google_cloud_run_service.lunch-bot.project
  service     = google_cloud_run_service.lunch-bot.name

  policy_data = data.google_iam_policy.lunch-bot-noauth.policy_data
}