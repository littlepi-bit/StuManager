#! /bin/bash

date

cd go
rm -rf new
mkdir new && cd new

git clone https://ghproxy.com/https://github.com/littlepi-bit/StuManager.git

cd StuManager

echo "Please input docker container ID: "
read x
docker stop $x
docker rm $x

echo "Please input docker images ID: "
read ys
docker rmi $y
echo "Please input image tag: "
read z
docker build -t stumanager:$z .
docker run --name stumanager -p 8000:8000 -d stumanager:$z

cd ..
cd ..
cd react 
rm -rf new
mkdir new && cd new

git clone https://ghproxy.com/https://github.com/littlepi-bit/StuManagerFront.git
cd StuManagerFront

echo "Please input docker container ID: "
read x
docker stop $x
docker rm $x

echo "Please input docker images ID: "
read y
docker rmi $y
echo "Please input image tag: "
read z
docker build -t stumanagerfront:$z .
docker run --name stumanagerfront -p 3000:3000 -d stumanagerfront:$z
docker ps -a