provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxxxxx"
  password = "xxxxxxx"
}

resource "bigip_command" "test-command" {
  name     = "command1"
  commands = ["show sys version"]
}