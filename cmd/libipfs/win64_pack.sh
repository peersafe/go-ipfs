echo "Win64 C lib start building..."
make build
echo "Win64 C lib end building..."

timestamp=`date +%Y%m%d`
win64name="libipfs.0.4.3_dev.win64."$timestamp".tar.gz"
echo $win32name

tar -zcvf $win64name libipfs.*

rm -rf libipfs.h
rm -rf libipfs.a

mv $win64name  /cygdrive/z/下载/