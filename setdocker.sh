#! /bin/bash

date
echo "If update StuManager?"
read x1
if [ $x1 = 'Y']
then
cd go
rm -rf new
mkdir new && cd new

git clone https://ghproxy.com/https://github.com/littlepi-bit/StuManager.git

cd StuManager
docker ps -a
echo "Please input docker container ID: "
read x
docker stop $x
docker rm $x
docker images
echo "Please input docker images ID: "
read y
docker rmi $y
echo "Please input image tag: "
read z
docker build -t stumanager:$z .
docker run --name stumanager -p 8000:8000 -d stumanager:$z
cd ..
cd ..
cd ..
fi
echo "if update StuFront?"
read x2
if [ $x2 = 'Y']
then
cd react
rm -rf new
mkdir new && cd new

git clone https://ghproxy.com/https://github.com/littlepi-bit/StuManagerFront.git
cd StuManagerFront
yarn
yarn build
echo "build again?"
read b
if [$b='Y']
then
yarn build
fi
cd build
vim index.html
docker ps -a
echo "Please input docker container ID: "
read x
docker stop $x
docker rm $x
docker images
echo "Please input docker images ID: "
read y
docker rmi $y
echo "Please input image tag: "
read z
docker build -t stufront:$z .
docker run --name stufront -p 3000:3000 -d stufront:$z
docker ps -a

docker images
echo "Please input docker image ID"
fi
