<h1>nagctl<h1>

a nagios cli control tool for nagios.cmd.
it can mute single or all hosts and services. will only take effect on hosts with OK status

<h5>pro:<h5>
This tool dose not require dependency what so ever as long as you are using bash!.

<h5>con:<h5>
result output part may not as good as other tool which can access json or print lib

<h5>before use:<h5>
change the location in script to wherever your nagios.cmd and status.dat are at


###############################################################################################
script will disable or enable notification for your hosts or services on nagios
if only services are pointed, script will use status.dat info to get all hosts 

Returns:
enable  -> 1
disable -> 0

current_state ->  0  => OK
current_state ->  1  => WARNING
current_state ->  2  => CRITICAL

example: nagctl -e enable -h hosts -s services
         nagctl -e diable -h hosts
         nagctl -e enable -s services  (this will only take effect on all hosts with OK status)
         								needs location of status.dat

notice: you can only use bash regex to discribe your hosts or services
		host -->> edge{12,13} or edge{12..20}
		service -->> {Logchute,http,write}  case sensitive
###############################################################################################

<h5>future plan:<h5>

more clear output

add check function to script
