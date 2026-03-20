# Basic PostgreSQL monitor with defaults
resource "bigip_gtm_monitor_postgresql" "basic" {
  name = "/Common/my_postgresql_monitor"
}

# PostgreSQL monitor with database authentication
resource "bigip_gtm_monitor_postgresql" "with_auth" {
  name                 = "/Common/my_pg_auth_monitor"
  defaults_from        = "/Common/postgresql"
  destination          = "*:5432"
  interval             = 10
  timeout              = 60
  probe_timeout        = 3
  ignore_down_response = "disabled"
  database             = "mydb"
  username             = "monitor_user"
  password             = "monitor_pass"
  receive              = "SELECT"
  debug                = "no"
}
