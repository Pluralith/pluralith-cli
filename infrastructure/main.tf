// - - Pluralith CLI CI/CD Infrastructure - -

terraform {
  backend "gcs" {
    bucket = "pluralith-infrastructure-state"
    prefix = "pluralith-cli"
  }
  
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}

provider "google" {
  project = var.project_id
}

// Google Cloud Storage
data "google_iam_policy" "pluralith_cli_bucket_policy" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
        "allUsers",
    ] 
  }
}

resource "google_storage_bucket_iam_policy" "pluralith_cli_bucket_policy_link" {
  bucket = google_storage_bucket.pluralith_cli_bucket.name
  policy_data = data.google_iam_policy.pluralith_cli_bucket_policy.policy_data
}

resource "google_storage_bucket" "pluralith_cli_bucket" {
  name = var.bucket_name
  location = var.bucket_location
}

// Google Cloud Build
resource "google_cloudbuild_trigger" "pluralith_api_cloud_build_trigger" {
  name = var.cloud_build_trigger_name
  description = var.cloud_build_trigger_description

  github {
    owner = var.repo_owner
    name = var.repo_name

    push {
      branch = var.repo_branch
    }
  }

  filename = "cloudbuild.yaml"
}