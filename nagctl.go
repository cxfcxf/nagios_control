package main

import (
	"fmt"
	"time"
	"strings"
	"regexp"
	"os/exec"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/auth"
	"github.com/vaughan0/go-ini"
	"github.com/cxfcxf/nagtomaps"
)

type DataPost struct {
	Exec 		string	`form:"exec"`
	Hosts 		string	`form:"hosts"`
	Services	string	`form:"services"`
}

type AckPost struct {
	Hosts  		string 	`form:"hosts"`
	Services   	string 	`form:"services"`
	Ackall		string 	`form:"ackall"`
}

type Dstatus struct {
	Servers_have_problem 	[]string
	Services_have_problem	map[string][]string
	Servers_enabled			[]string
	Services_enabled		map[string][]string
	Servers_ok_disabled		[]string
	Services_ok_disabled 	map[string][]string
}

func iniParser(inifile string) (string, string, string, string, string) {
	f, err := ini.LoadFile(inifile)
	if err != nil {panic(err)}
	authname, _ := f.Get("authlogin", "username")
	authpass, _ := f.Get("authlogin", "password")
	port, _ := f.Get("listen", "port")
	sfile, _ := f.Get("files", "sfile")
	efile, _ := f.Get("files", "efile")
	return authname, authpass, port, sfile, efile
}

func getStatus(sfile string) Dstatus{
	sdata := nagtomaps.ParseStatus(sfile)

	var dstatus Dstatus

	//hosts have problem
	for _, server := range sdata.Hoststatuslist {
		if  server["current_state"] != "0"  && server["notifications_enabled"] == "1" && server["acknowledgement_type"] == "0" && server["current_attempt"] > "1" {
			dstatus.Servers_have_problem = append(dstatus.Servers_have_problem, server["host_name"])
		}
	}

	//hosts enabled
	for _, server := range sdata.Hoststatuslist {
		if  (server["acknowledgement_type"] != "0" && server["current_state"] != "0") || (server["notifications_enabled"] == "0" && server["current_state"] != "0") {
		} else {
			dstatus.Servers_enabled = append(dstatus.Servers_enabled, server["host_name"])
		}
	}

	//host ok disabled
	for _, server := range sdata.Hoststatuslist {
		if  server["current_state"] == "0" && server["notifications_enabled"] == "0" {
			dstatus.Servers_ok_disabled = append(dstatus.Servers_ok_disabled, server["host_name"])
		}
	}

	//services have problem
	dstatus.Services_have_problem = make(map[string][]string)

	for _, serverserv := range sdata.Servicestatuslist {
		for _, service := range serverserv {
			if service["current_state"] != "0" && service["notifications_enabled"] == "1" && service["acknowledgement_type"] == "0" && service["current_attempt"] > "1" {
				dstatus.Services_have_problem[service["host_name"]] = append(dstatus.Services_have_problem[service["host_name"]], service["service_description"])
			}
		}
	}

	//services enabled
	dstatus.Services_enabled = make(map[string][]string)

	for _, serverserv := range sdata.Servicestatuslist {
		for _, service := range serverserv {
			if (service["acknowledgement_type"] != "0" && service["current_state"] != "0") || (service["notifications_enabled"] == "0" && service["current_state"] != "0") {
			} else {
				dstatus.Services_enabled[service["host_name"]] = append(dstatus.Services_enabled[service["host_name"]], service["service_description"])
			}
		}
	}

	//services ok disabled
	dstatus.Services_ok_disabled = make(map[string][]string)

	for _, serverserv := range sdata.Servicestatuslist {
		for _, service := range serverserv {
			if service["current_state"] == "0" && service["notifications_enabled"] == "0" {
				dstatus.Services_ok_disabled[service["host_name"]] = append(dstatus.Services_ok_disabled[service["host_name"]], service["service_description"])
			}
		}
	}

	return dstatus
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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

func nagiosExecCtl(execute string, hosts string, services string, ds Dstatus, efile string) {
	exe := strings.ToUpper(execute)

	if hosts != "" && services != "" {
		h, _ := regexp.Compile(hosts)
		s, _ := regexp.Compile(services)

		for host, secs := range ds.Services_enabled {
			for _, service := range secs {
				if h.MatchString(host) && s.MatchString(service) {
					command := fmt.Sprintf("/bin/echo \"[%d] %s_SVC_NOTIFICATIONS;%s;%s\n\" > %s", time.Now().Unix(), exe, host, service, efile)
					cmd := exec.Command("sh", "-c", command)
					if err := cmd.Run() ; err != nil {panic(err)}
				}
			}
		}
	} else if hosts != "" && services == "" {
		h, _ := regexp.Compile(hosts)

		for _, host := range ds.Servers_enabled {
			if h.MatchString(host) {
				command := fmt.Sprintf("/bin/echo \"[%d] %s_HOST_NOTIFICATIONS;%s\n\" > %s", time.Now().Unix(), exe, host, efile)
				cmd := exec.Command("sh", "-c", command)
				if err := cmd.Run() ; err != nil {panic(err)}
				command = fmt.Sprintf("/bin/echo \"[%d] %s_HOST_SVC_NOTIFICATIONS;%s\n\" > %s", time.Now().Unix(), exe, host, efile)
				cmd = exec.Command("sh", "-c", command)
				if err := cmd.Run() ; err != nil {panic(err)}
			}
		}
	} else if hosts == "" && services != "" {
		s, _ := regexp.Compile(services)

		for host, secs := range ds.Services_enabled {
			for _, service := range secs {
				if s.MatchString(service) {
					command := fmt.Sprintf("/bin/echo \"[%d] %s_SVC_NOTIFICATIONS;%s;%s\n\" > %s", time.Now().Unix(), exe, host, service, efile)
					cmd := exec.Command("sh", "-c", command)
					if err := cmd.Run() ; err != nil {panic(err)}
				}
			}
		}
	} else if hosts == "" && services == "" {
		for _, host := range ds.Servers_ok_disabled {
				command := fmt.Sprintf("/bin/echo \"[%d] %s_HOST_NOTIFICATIONS;%s\n\" > %s", time.Now().Unix(), exe, host, efile)
				cmd := exec.Command("sh", "-c", command)
				if err := cmd.Run() ; err != nil {panic(err)}
		}

		for host, secs := range ds.Services_ok_disabled {
			for _, service := range secs {
				command := fmt.Sprintf("/bin/echo \"[%d] %s_SVC_NOTIFICATIONS;%s;%s\n\" > %s", time.Now().Unix(), exe, host, service, efile)
				cmd := exec.Command("sh", "-c", command)
				if err := cmd.Run() ; err != nil {panic(err)}
			}
		}
	}
}

func main() {
	authname, authpass, port, sfile, efile := iniParser("nagctl.ini")

	m := martini.Classic()

	m.Use(render.Renderer())

	m.Use(auth.Basic(authname, authpass))

	m.Get("/", func() string {
		return "you are hiting the webroot directory! please use /status or /nagctl"	
	})

	m.Get("/status", func(r render.Render) {
		ds := getStatus(sfile)
		r.HTML(200, "status", ds)
	})

	m.Post("/status", binding.Bind(AckPost{}), func(ap AckPost, r render.Render) {
		if ap.Ackall == "notempty" {
			ds := getStatus(sfile)
			for _, s := range ds.Servers_have_problem {
				go nagiosExec(s, "", efile)
			}
			for sr, svs := range ds.Services_have_problem {
				if !stringInSlice(sr, ds.Servers_have_problem) {
					for _, sev := range svs {
						go nagiosExec(sr, sev, efile)
					}
				}
			}
			ap.Ackall = fmt.Sprintf("Servers: ==> %s\n Services: ==> %s\n", ds.Servers_have_problem, ds.Services_have_problem)
		} else if ap.Hosts != "" {
			go nagiosExec(ap.Hosts, "", efile)
		} else {
			hslist := strings.Split(ap.Services, " ")
			go nagiosExec(hslist[0], hslist[1], efile)
		}
		r.HTML(200, "finish", ap)
	})

	m.Get("/nagctl", func(r render.Render) {
		r.HTML(200, "nagctl", nil)
	})

	m.Post("/nagctl", binding.Bind(DataPost{}), func(dp DataPost, r render.Render) {
		ds := getStatus(sfile)
		nagiosExecCtl(dp.Exec, dp.Hosts, dp.Services, ds, efile)
		command := fmt.Sprintf("Executed --> Exec: %s  Hosts: %s  services: %s \n",
			dp.Exec, dp.Hosts, dp.Services)
		r.HTML(200, "nagfin", command)
	})

	m.RunOnAddr(port)
}
