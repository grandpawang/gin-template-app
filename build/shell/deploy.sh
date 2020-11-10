pid=$(ssh root@192.168.31.51 " lsof -i:50052 | grep TCP | awk '{print $2}' ")
if [ -n "$pid" ]; then
    echo "kill :50052 program"
    ssh root@192.168.31.51 "kill $pid"
fi

set -x &&
    ssh root@192.168.31.51 "rm -rf /home/coint/gbbmn" &&
    ssh root@192.168.31.51 "mkdir /home/coint/gbbmn && chmod 777 /home/coint/gbbmn" &&
    cd cmd &&
    scp favicon.ico root@192.168.31.51:/home/coint/gbbmn/favicon.ico &&
    scp cloud root@192.168.31.51:/home/coint/gbbmn/cloud &&
    scp config.toml root@192.168.31.51:/home/coint/gbbmn/config.toml &&
    ssh root@192.168.31.51 "nohup /home/coint/gbbmn/cloud --conf=/home/coint/gbbmn/config.toml --icon=/home/coint/gbbmn/favicon.ico >/home/coint/gbbmn/log 2>&1 &" &&
    ls
