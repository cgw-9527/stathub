#!/bin/bash
echo "start--------------"
STATHUB_URL="http://code.tianhecloud.com:33012/chenguowen/stathub/src/master/"

BASEDIR="/usr/local/stathub"

$sudo mkdir -p $BASEDIR
echo "after mkdir"
$sudo chown -R $(id -u -n):$(id -g -n) $BASEDIR
echo "after chwon"
if [ ! -d $BASEDIR ]; then
    echo "Unable to create dir $BASEDIR and chown to current user, Please manual do it"
    exit 1
fi
cd $BASEDIR

for i in "i686"; do
    if command_exists wget; then
        wget --no-check-certificate "$STATHUB_URL/stathub.$i.zip"
    elif command_exists curl; then
        curl --insecure -O "$STATHUB_URL/stathub.$i.zip"
    else
        echo "Unable to find curl or wget, Please install it and try again."
        exit 1
    fi
done

if [ ! -f stathub.i686.zip ]; then
    echo "Unable to get server file, Please manual download it"
    exit 1
fi

unzip stathub.i686.zip
chmod +x stathub service
[ ! -d conf ] && $sudo mkdir $BASEDIR/conf
if [ ! -f conf/stathub.conf ]; then
    $sudo ./stathub -c conf/stathub.conf --init-server
fi
$sudo mkdir $BASEDIR/pkgs
mv stathub.*.zip $BASEDIR/pkgs

if [ -z "$(grep stathub /etc/rc.local)" ]; then
    $sudo sh -c "echo \"cd $BASEDIR;sudo rm -f log/stathub.pid; ./service start >>$BASEDIR/log/stathub.log 2>&1\" >> /etc/rc.local"
fi

echo "----------------------------------------------------"
echo "| Server install successful, Please start it using |"
echo "| ./service {start|stop|restart}                   |"
echo "| Now it will automatic start                      |"
echo "|                                                  |"
echo "| Feedback: https://github.com/likexian/stathub-go |"
echo "| Thank you for your using, By Li Kexian           |"
echo "| StatHub, Apache License, Version 2.0             |"
echo "----------------------------------------------------"

$sudo ./service start
