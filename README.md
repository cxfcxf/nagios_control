# nagctl

## This is a nagios control web interface

unlike the original nagios webui, nagctl.go can be used as an addon for nagios

```
usage:
	go run nagctl.go -sfile status.dat -efile nagios.cmd
```

### this will start a web server on your port :3333 , default login is username. password.

it has two sub directory

```
/status
for quick acknowledge alert (mobile device compatiable

/nagctl
you can use regex to point the server or service you would like to mute/unmute alert for
```

### future

would put all those config into a ini file. and you set stuff there.

if you have any other suggestion please let me know
siegfried.chen@gmail.com
