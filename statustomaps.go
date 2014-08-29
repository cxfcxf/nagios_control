package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func loadStatus(statusfile string, s string) []string{
	statuses := make([]string, 0)
	var statusfile_temp string
	for {
		nextOpenBracket := strings.Index(statusfile, s)

		if nextOpenBracket == -1 { break }

		statusfile_temp = statusfile[nextOpenBracket:]

		nextCloseBracket := strings.Index(statusfile_temp, "}") + nextOpenBracket
		statuses = append(statuses, statusfile[nextOpenBracket:nextCloseBracket+1])
		statusfile = statusfile[nextCloseBracket:]
	}
	return statuses
}

func infoStatus(infostatus []string) map[string]string{
	infostatuslist := make(map[string]string)

	for _, is := range infostatus {
		isline := strings.Split(is, "\n")
		isline = isline[1:len(isline)-1]

		for _, isl := range isline {
			line := strings.Split(isl, "=")
			infostatuslist[strings.TrimSpace(line[0])] = string(line[1])
		}
	}
	return infostatuslist
}

func hostStatus(hoststatuses []string) map[string]map[string]string{
	hoststatuslist := make(map[string]map[string]string)

	for _, hs := range hoststatuses {
		hsline := strings.Split(hs, "\n")
		hsline = hsline[1:len(hsline)-1]

		hsname := strings.Split(string(hsline[0]), "=")[1]

		if _, ok := hoststatuslist[hsname]; !ok {
			hoststatuslist[hsname] = make(map[string]string)
		}

		for _, hsl := range hsline {
			line := strings.Split(hsl, "=")
			hoststatuslist[hsname][strings.TrimSpace(line[0])] = string(line[1])
		}
	}
	return hoststatuslist
}

func serviceStatus(servicestatuses []string) map[string]map[string]map[string]string{
	
	servicestatuslist  := make(map[string]map[string]map[string]string)

	for _, ss := range servicestatuses {
		ssline := strings.Split(ss, "\n")
		ssline = ssline[1:len(ssline)-1]

		hsname := strings.Split(string(ssline[0]), "=")[1]
		servdesc := strings.Split(string(ssline[1]), "=")[1]

		if _, ok := servicestatuslist[hsname]; !ok {
			servicestatuslist[hsname] = make(map[string]map[string]string)
		}
		if _, ok := servicestatuslist[hsname][servdesc]; !ok {
			servicestatuslist[hsname][servdesc] = make(map[string]string)
		}

		for _, ssl := range ssline {
			line := strings.Split(ssl, "=")
			servicestatuslist[hsname][servdesc][strings.TrimSpace(line[0])] = string(line[1])
		}
	}
	return servicestatuslist
}

type statusData struct {
	Infostatuslist		map[string]string
	Programstatuslist	map[string]string
	Hoststatuslist		map[string]map[string]string
	Servicestatuslist	map[string]map[string]map[string]string
	Hostcommentslist	map[string]map[string]string
	Servicecommentslist	map[string]map[string]string
}

func parseStatus(nagstatfile string)  statusData{
	b, err := ioutil.ReadFile(nagstatfile)
	if err != nil { panic(err) }

	var sdata statusData

	infostatuses := loadStatus(string(b), "info {")
	sdata.Infostatuslist = infoStatus(infostatuses)

	programstatuses := loadStatus(string(b), "programstatus {")
	sdata.Programstatuslist = infoStatus(programstatuses)

	hoststatuses := loadStatus(string(b), "hoststatus {")
	sdata.Hoststatuslist = hostStatus(hoststatuses)

	servicestatuses := loadStatus(string(b), "servicestatus {")
	sdata.Servicestatuslist = serviceStatus(servicestatuses)

	hostcomments := loadStatus(string(b), "hostcomment {")
	sdata.Hostcommentslist = hostStatus(hostcomments)

	servicecomments := loadStatus(string(b), "servicecomment {")
	sdata.Servicecommentslist = hostStatus(servicecomments)

	return sdata
}








