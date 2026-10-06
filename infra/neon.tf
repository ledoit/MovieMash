# Free-plan Postgres. The provider uses NEON_API_KEY.
# Default project settings stay on Neon's free plan; do not raise compute past 0.25 CU.
resource "neon_project" "moviemash" {
  name       = "moviemash"
  pg_version = 16

  primary_compute {
    autoscaling_limit_min_cu = 0.25
    autoscaling_limit_max_cu = 0.25
    suspend_timeout_seconds  = 300
  }
}

output "database_url" {
  description = "Postgres URL for the API. Sensitive."
  value       = neon_project.moviemash.connection_uri
  sensitive   = true
}
