#!/bin/bash
cd ~
rm -rf ulord
mkdir ulord
cd ulord
wget ftp://tools.ulord.one/ulord_1_1_86.tgz
tar zxf ulord_1_1_86.tgz
mkdir success
rm -rf /usr/local/bin/ulord-cli
rm -rf /usr/local/bin/ulordd
cp -rf ulordd ulord-cli success /usr/local/bin
