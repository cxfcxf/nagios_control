package main

import (
	"fmt"
	"encoding/json"
	"github.com/cxfcxf/nagtomaps"
)

func main() {
	
	sdata := nagtomaps.ParseStatus("status.dat")
	// status.dat is the file from nagios
	fmt.Println(sdata.Infostatuslist)
	//print maps

	j, _ := json.Marshal(sdata.Hoststatuslist)
	fmt.Println(string(j))
	//print json

	fmt.Println(sdata.Servicestatuslist["edge21-fra"])
	// print map of all service description
	fmt.Println(sdata.Servicestatuslist["edge21-fra"]["HTTP"]["notifications_enabled"])
	// print status like 1 or 0
	
}
