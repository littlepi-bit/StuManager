#! /bin/bash

date

rm -rf new
mkdir new && cd new

git clone https://ghproxy.com/https://github.com/littlepi-bit/StuManager.git

cd StuManager

echo "Please input docker container ID: "
read x
docker stop $x
docke rm $x

echo "Please input docker images ID: "
read y
docker rmi $y
echo "Please input image tag: "
read z
docker build -t stumanager:$z .
docker run --name stumanager -p 8000:8000 -d stumanager:$z