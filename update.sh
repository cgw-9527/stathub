#!/bin/bash
cd ~
rm -rf ulord
mkdir ulord
cd ulord
wget ftp://tools.ulord.one/ulord_1_1_86.tgz
tar -zxf ulord_1_1_86.tgz
rm -rf /usr/local/bin/ulord-cli
rm -rf /usr/local/bin/ulordd
cp -i ulordd ulord-cli /usr/local/bin