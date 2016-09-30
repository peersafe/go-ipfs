#!/bin/sh

echo "ios framework start building......"
make ios
echo "ios framework end building......"

echo "android framework start building......"
make android
echo "android framework end building......"


#timestamp=`date +%Y%m%d%H%M%S` 
timestamp=`date +%Y%m%d` 
iosname="libipfs.0.4.3_dev.ios."$timestamp".tar.gz"
androidname="libipfs.0.4.3_dev.android."$timestamp".tar.gz"
echo $iosname
echo $androidname

echo "pack ios..."
tar -zcvf $iosname Ipfsmobile.framework

echo "pack android..."
tar -zcvf $androidname ipfsmobile.aar

rm -rf Ipfsmobile.framework
rm -rf ipfsmobile.aar
