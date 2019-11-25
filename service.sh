#!/bin/bash
PIDFILE="log/stathub.pid"

start() {
        echo "starting"
        sudo nohup ./stathub -c conf/stathub.conf & \Entsdaer
}

stop() {
    echo "stopping"
    if [ -f $PIDFILE ]; then
        kill -9 `cat $PIDFILE`
        rm -rf $PIDFILE
        echo "ok"
    fi
}

case "$1" in
    start)
        start ;;
    stop)
        stop
        ;;
    restart)
        stop
        start
        ;;
    *)
        echo "Usage: $0 {start|stop|restart}"
        exit 1
esac
