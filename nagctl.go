package main

import (
	"fmt"
	"os/exec"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/binding"
)

var flag (
	-s nagios stats file
	-c nagios cmd file
)


type DataPost struct {
	Exec	string	`form:"exec"`
	Hosts	string	`form:"hosts"`
	Servers	string	`form:"services"`
}

func nagiosState(sf string) string{
	status, err := ioutil.Readfile(sf)
	if err != nil {panic(err)}
	return status
}

func nagiosExec(e string, h string, s string) string {
	for _, e := range e {
		go exec(e)
	}
}

func main() {
	flag.Parse()

	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func() string {
		return "you are hiting the webroot directory! please use /nagctl"	
	})

	m.Get("/nagctl", func(r render.Render) {
		status := nagiosState(*s)
		r.HTML(200, "nagctl", status)
	})

	m.Post("/nagctl", binding.Bind(DataPost{}), func(dp DataPost) string {
		res := nagiosExec(dp.Exec, dp.Hosts, dp.Servers)
		return fmt.Sprintf("Executed --> Exec: %s  Hosts: %s  services: %s \n output: \n %s",
			dp.Exec, dp.Hosts, dp.Servers, res)

	})

	m.Run()
}