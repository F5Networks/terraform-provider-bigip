data "bigip_ltm_pool" "test" {
  name      = "terraform-pool"
  partition = "Common"
}

output "test" {
  value = data.bigip_ltm_pool.test.full_path
}