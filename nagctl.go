package main

import (
	"fmt"
	"flag"
	"time"
	"strings"
	"os/exec"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/binding"
	"github.com/cxfcxf/nagtomaps"
)

var (
	sfile = flag.String("sfile", "", "path to the status.dat file nagios create")
	efile = flag.String("efile", "", "path to nagios external command file nagios.cmd")
)

type DataPost struct {
	Exec	string	`form:"exec"`
	Hosts	string	`form:"hosts"`
	Servers	string	`form:"services"`
}

type AckPost struct {
	Hosts  		string 	`form:"hosts"`
	Services   	string 	`form:"services"`
	Ackall		string 	`form:"ackall"`
}

type Dstatus struct {
	Servers_have_problem 	[]string
	Services_have_problem	map[string][]string
}

func getStatus(sfile string) Dstatus{
	sdata := nagtomaps.ParseStatus(sfile)

	var dstatus Dstatus

	//hosts
	for _, server := range sdata.Hoststatuslist {
		if  server["current_state"] != "0"  && server["notifications_enabled"] == "1" && server["acknowledgement_type"] == "0" {
			dstatus.Servers_have_problem = append(dstatus.Servers_have_problem, server["host_name"])
		}
	}
	//services
	dstatus.Services_have_problem = make(map[string][]string)

		for _, serverserv := range sdata.Servicestatuslist {
			for _, service := range serverserv {
				if service["current_state"] != "0" && service["notifications_enabled"] == "1" && service["acknowledgement_type"] == "0" {
					dstatus.Services_have_problem[service["host_name"]] = append(dstatus.Services_have_problem[service["host_name"]], service["service_description"])
				}
			}
		}

	return dstatus
}

func nagiosExec(host string, service string,  efile string) {
	if host != "" && service == "" {
		command := fmt.Sprintf("/bin/echo \"[%d] ACKNOWLEDGE_HOST_PROBLEM;%s;2;1;0;admin;acknowledged by nagctl webui\n\" > %s", time.Now().Unix(), host, efile)
		cmd := exec.Command("sh", "-c", command)
		if err := cmd.Run() ; err != nil {panic(err)}
	} else {
		command := fmt.Sprintf("/bin/echo \"[%d] ACKNOWLEDGE_SVC_PROBLEM;%s;%s;2;1;0;admin;acknowledged by nagctl webui\n\" > %s", time.Now().Unix(), host, service, efile)
		cmd := exec.Command("sh", "-c", command)
		if err := cmd.Run() ; err != nil {panic(err)}
	}
}

func main() {
	flag.Parse()

	m := martini.Classic()

	m.Use(render.Renderer())

	m.Get("/", func() string {
		return "you are hiting the webroot directory! please use /nagctl"	
	})

	m.Get("/status", func(r render.Render) {
		ds := getStatus(*sfile)
		r.HTML(200, "status", ds)
	})

	m.Post("/status", binding.Bind(AckPost{}), func(ap AckPost, r render.Render) {
		if ap.Ackall != "" {
			ds := getStatus(*sfile)
			for _, s := range ds.Servers_have_problem {
				go nagiosExec(s, "", *efile)
			}
			for sr, svs := range ds.Services_have_problem {
				for _, sev := range svs {
					go nagiosExec(sr, sev, *efile)
				}
			}
		} else if ap.Hosts != "" {
			go nagiosExec(ap.Hosts, "", *efile)
		} else {
			hslist := strings.Split(ap.Services, " ")
			go nagiosExec(hslist[0], hslist[1], *efile)
		}
		r.HTML(200, "finish", ap)
	})

	m.Get("/nagctl", func(r render.Render) {
		r.HTML(200, "nagctl", nil)
	})

	m.Post("/nagctl", binding.Bind(DataPost{}), func(dp DataPost) string {
		return fmt.Sprintf("Executed --> Exec: %s  Hosts: %s  services: %s \n output: \n %s",
			dp.Exec, dp.Hosts, dp.Servers)
	})

	m.RunOnAddr(":3333")
}