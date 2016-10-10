echo "Win32 C lib start building..."
make build
echo "Win32 C lib end building..."

timestamp=`date +%Y%m%d`
win32name="libipfs.0.4.3_dev.win32."$timestamp".tar.gz"
echo $win32name

tar -zcvf $win32name libipfs.*

rm -rf libipfs.h
rm -rf libipfs.a

mv $win32name  /cygdrive/z/下载/