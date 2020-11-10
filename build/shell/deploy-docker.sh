ssh root@192.168.31.51 "rm -rf /home/coint/gbbmn/gbbmn-cloud"
ssh root@192.168.31.51 "mkdir -p /home/coint/gbbmn/gbbmn-cloud/ && chmod 777 /home/coint/gbbmn/gbbmn-cloud/"
scp -r cmd root@192.168.31.51:/home/coint/gbbmn/gbbmn-cloud/
scp dockerfile root@192.168.31.51:/home/coint/gbbmn/gbbmn-cloud/
ssh root@192.168.31.51 '''
exist=$(docker ps -a | grep gbbmn-cloud)
if [ -n "$exist" ]; then
    echo exist gbbmn-cloud docker
    docker stop gbbmn-cloud
    docker rm gbbmn-cloud
    docker rmi gbbmn-cloud
fi
'''
ssh root@192.168.31.51 '''
set -x &&
    docker build -t gbbmn-cloud /home/coint/gbbmn/gbbmn-cloud/ &&
    docker run -d -p 50052:50052 --name=gbbmn-cloud gbbmn-cloud
'''
