# nagctl

## This is a nagios control web interface

unlike the original nagios webui, nagctl.go can be used as an addon for nagios

```
usage:
	go run nagctl.go
```

### update iniParser, setting is now on nagctl.ini. its pretty self-explain

it has two sub directory

```
/status
for quick acknowledge alert (mobile device compatiable)

/nagctl
you can use regex to point the server or service you would like to mute/unmute alert for

if you submit without host and service, keep exec enable it will clean false mute
```

### future

add more functions and create a dashboard for desktop browser

if you have any other suggestion please let me know
siegfried.chen@gmail.com
