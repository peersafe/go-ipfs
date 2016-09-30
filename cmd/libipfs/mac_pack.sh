echo "mac C lib start building..."
make build
echo "Mac C lib end building..."

timestamp=`date +%Y%m%d%H%M%S`
macname="libipfs.0.4.3_dev.Mac."$timestamp".tar.gz"
echo $macname

tar -zcvf $macname libipfs.*

rm -rf libipfs.h
rm -rf libipfs.a

mv $macname /Users/sunzhiming/Downloads
