# Free-plan Postgres. The provider uses NEON_API_KEY.
# Default project settings stay on Neon's free plan; do not raise compute past 0.25 CU.
resource "neon_project" "moviemash" {
  name       = "moviemash"
  org_id     = "org-odd-snow-78672228"
  pg_version = 16
  # Free-plan ceiling. The provider default of 86400 is rejected.
  history_retention_seconds = 21600

  default_endpoint_settings {
    autoscaling_limit_min_cu = 0.25
    autoscaling_limit_max_cu = 0.25
  }
}

output "database_url" {
  description = "Postgres URL for the API. Sensitive."
  value       = neon_project.moviemash.connection_uri
  sensitive   = true
}
