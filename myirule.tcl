when HTTP_REQUEST {
 
  if { [string tolower [HTTP::header value Upgrade]] equals "websocket" } {
    HTTP::disable
#    ASM::disable
    log local0. "[IP::client_addr] - Connection upgraded to websocket protocol. Disabling ASM-checks and HTTP protocol. Traffic is treated as L4 TCP stream."
  } else {
    HTTP::enable
#    ASM::enable
    log local0. "[IP::client_addr] - Regular HTTP request. ASM-checks and HTTP protocol enabled. Traffic is deep-inspected at L7."
  }
}
 
