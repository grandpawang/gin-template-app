#!/bin/bash
# 设置mysql的登录用户名和密码(根据实际情况填写)
mysql_user="root"
mysql_password="1234"
mysql_host="localhost"
mysql_port="3307"
#mysql_charset="utf8mb4"

# 备份文件存放地址(根据实际情况填写)
#backup_location=/var/lib/docker/backup
backup_location=/home/db_backup

backup_time=$(date +%Y%m%d%H%M)
# shellcheck disable=SC2164
cd $backup_location
DumpFile=gbbmn-$backup_time.sql
GZDumpFile=gbbmn-$backup_time.sql.tgz

if [ ! -d $backup_location ];  then
	mkdir -p $backup_location
fi

# 是否删除过期数据
expire_backup_delete="ON"
expire_days=7


#welcome_msg="Welcome to use MySQL backup tools!"
# 备份指定数据库中数据(此处设定数据库是gbbmn)
# docker的mysql名称：mysql-master
docker_mysql=mysql-master

docker exec $docker_mysql mysqldump -h$mysql_host -P$mysql_port -u$mysql_user -p$mysql_password -B gbbmn > $DumpFile

# 打包sql文件压缩
/bin/tar -czvf $GZDumpFile $DumpFile
/bin/rm $DumpFile

# 删除过期数据
# shellcheck disable=SC2166
if [ "$expire_backup_delete" == "ON" -a  "$backup_location" != "" ];then
        # shellcheck disable=SC2038
        # shellcheck disable=SC2091
        $(find $backup_location -type f -mtime +$expire_days | xargs rm -rf)
        echo "Expired backup data delete complete!"
fi


# 新建定时任务

##新建定时任务命令
#crontab -e
##查看定时任务命令
#crontab -l
##删除所有定时任务命令
#crontab -r

# 0 1 * * 1-6 cd /root/docker;sh backupDocker.sh >> /dev/docker/log.txt 2>>/dev/docker/log.txt
# 0 17 * * 0-6 cd /home;sh backup.sh >> /home/log.txt 2>&1


# 恢复数据
# 进入Docker-mysql
# use gbbmn;
# source /var/lib/docker/backup/gbbmn-xxx.sql

