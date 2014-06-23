#!/bin/bash

#by siegfried chen

#nagios contrl script
#for now notification only

usage()
{
cat << EOF
##########################################################################################
script will disable or enable notification for your hosts or services on nagios
if only services are pointed, script will use status.dat info to get all hosts 

Returns:
enable  -> 1
disable -> 0

current_state ->  0  => OK
current_state ->  1  => WARNING
current_state ->  2  => CRITICAL

example: nagctl -e list -h hosts
         nagctl -e enable -h hosts -s services
         nagctl -e diable -h hosts
         nagctl -e enable -s services  (this will take effect on all hosts with OK status)


notice: you can only use bash regex to discribe your hosts or services
		host -->> "edge{12,13}" or "edge{12..20}"
		service -->> "{LOGCHUTE,WRITES}"  case sensitive
##########################################################################################
EOF
}

listall() 
{
H=`eval echo $1`
for x in $H; do
echo "HOST===>${x}"
awk -v RS= -v ORS="\n\n" "/${x}/" ${STATFILE}|grep -E 'service_desc|current_state|notifications_enabled'|awk -F= '{print $1, "--->", $2}'
done
}


while getopts “e:h:s:” OPTION
do
     case $OPTION in
         e)
             EXEC=$OPTARG
             ;;
         h)
             HOST=$OPTARG
             ;;
         s)
             SERV=$OPTARG
             ;;
         ?)
             usage
             exit
             ;;
     esac
done

CMD="/usr/local/nagios/var/rw/nagios.cmd"
NOW=`date +%s`
STATFILE="/usr/local/nagios/var/status.dat"


if [ "$EXEC" == "list" ];then
	listall $HOST
	exit 0
fi


if [ -z $EXEC ] || ([ -z $HOST ] && [ -z $SERV ]); then
     usage
     exit 1
fi

if [ "$EXEC" == "enable" ] || [ "$EXEC" == "disable" ]; then
	EXEC=`echo $EXEC|awk '{print toupper($0)}'`
	if [ -n $HOST ] && [ -z $SERV ]; then
		HOST=`eval echo $HOST`
		for x in $HOST; do
			echo "[${NOW}] ${EXEC}_HOST_NOTIFICATIONS;${x}" > $CMD
			echo "[${NOW}] ${EXEC}_HOST_SVC_NOTIFICATIONS;${x}" > $CMD
		done
	elif [ -n "$HOST" ] && [ -n "$SERV" ]; then
		HOST=`eval echo $HOST`
		SERV=`eval echo $SERV`
		for x in $HOST; do
			for y in $SERV; do
				echo "[${NOW}] ${EXEC}_SVC_NOTIFICATIONS;${x};${y}" > $CMD
			done
		done
	else   #apply servises to site wide exlude hosts with state WARNING or CRITICAL
		ALLOT=(`awk -v RS= -v ORS="\n\n" '/current_state=(1|2)/' $STATFILE|grep "host_name" |awk -F= '{print $2}'|sort |uniq`)
		ALLHOST=(`awk -v RS= -v ORS="\n\n" '/current_state=0/' $STATFILE|grep "host_name" |awk -F= '{print $2}'|sort| uniq`)
		ALLZ=()

		for x in ${ALLHOST[@]};do
			SWITCH=0
			for y in ${ALLOT[@]};do
				if [ "$x" == "$y" ];then
					SWITCH=1
				fi
			done
			if [ "$SWITCH" -eq 0 ];then
				ALLZ+=("$x")
			fi
		done

		SERV=`eval echo $SERV`
		for x in ${ALLZ[@]}; do
			for y in $SERV; do
				echo "[${NOW}] ${EXEC}_SVC_NOTIFICATIONS;${x};${y}" > $CMD
			done
		done
		echo "This will not return any thing! go check nagios web"
		exit 0
	fi
fi


sleep 10
listall "$HOST"


