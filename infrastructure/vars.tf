variable "project_id" {
  type = string
  description = "GCP Project ID"
}

variable "bucket_name" {
  type = string
  description = "Pluralith CLI artifact bucket name"
}

variable "bucket_location" {
  type = string
  description = "Pluralith CLI artifact bucket location"
}

variable "cloud_build_trigger_name" {
  type = string
  description = "Cloud Build Trigger Name"
}

variable "cloud_build_trigger_description" {
  type = string
  description = "Cloud Build Trigger Description"
}

variable "repo_owner" {
  type = string
  description = "Github Repository Owner"
}

variable "repo_name" {
  type = string
  description = "Github Repository Name"
}

variable "repo_branch" {
  type = string
  description = "Github Repository Branch"
}