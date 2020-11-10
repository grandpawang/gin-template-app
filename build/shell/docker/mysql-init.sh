#!/bin/bash
# ssh root@192.168.31.51

function masterAccess() {
  docker exec -it mysql-master bash -c "
mysql -u root -p1234 <<EOF
# 主库创建同步用户 主库给同步用户授权
use mysql;
CREATE USER 'repl'@'%' IDENTIFIED WITH mysql_native_password BY 'repl@123456';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%';
FLUSH PRIVILEGES;
# 查询主库状态，并记录 File 的值和 Position 的值
SHOW MASTER STATUS;
EOF
"
  echo "🎃 mysql master 🚗 exit"
}

function slave1copy() {
  echo "🎈 mysql slave1 🚑"
  docker exec -it mysql-slave1 bash -c "
mysql -u root -p1234 <<EOF
change master to
  master_host='192.168.31.51',
  master_user='repl',
  master_log_file='mysql-bin.000003',
  master_log_pos=$position,
  master_port=3307,
  master_password='repl@123456';
start slave;
show slave status \G;
EOF
"
  echo "🎃 mysql slave1 🚑 exit"
}

function slave2copy() {
  echo "🎈 mysql slave2 🚍"
  docker exec -it mysql-slave2 bash -c "
mysql -u root -p1234 <<EOF
change master to
  master_host='192.168.31.51',
  master_user='repl',
  master_log_file='mysql-bin.000003',
  master_log_pos=$position,
  master_port=3307,
  master_password='repl@123456';
start slave;
show slave status \G;
EOF
"
  echo "🎃 mysql slave2 exit 🚍"
}

echo "🎈 mysql master 🚗"
position=$(docker exec -it mysql-master bash -c "mysql -u root -p1234 -e \"SHOW MASTER STATUS;\" " | grep mysql-bin.000003 | awk '{print $4}' | xargs echo)
create=$(docker exec -it mysql-master bash -c "mysql -u root -p1234 -e \"use mysql; select count(1) c from user where user='repl';\" " | grep 0)
echo position: $position
echo create: $create
if [[ "$create" != "" ]]; then
  masterAccess
  slave1copy
  slave2copy
  docker exec -it mysql-master bash -c "mysql -u root -p1234 -e \"create database gbbmn DEFAULT CHARACTER SET utf8mb4;\" "
  docker exec -it mysql-master bash -c "mysql -u root -p1234 -e \"create database grafana DEFAULT CHARACTER SET utf8mb4;;\" "

else
  # 查看主从复制
  slave1_IO=$(docker exec -it mysql-slave1 bash -c "mysql -u root -p1234 -e \"show slave status \G;\" " | grep Slave_IO_Running: | awk '{print $2}')
  slave1_SQL=$(docker exec -it mysql-slave1 bash -c "mysql -u root -p1234 -e \"show slave status \G;\" " | grep Slave_SQL_Running: | awk '{print $2}')
  echo "Slave_IO_Running: ${#slave1_IO} $slave1_IO"
  echo "Slave_SQL_Running: ${#slave1_SQL} $slave1_SQL"

  if ((${#slave1_IO} < 3 && ${#slave1_SQL} < 3)); then
    slave1copy
  fi

  slave2_IO=$(docker exec -it mysql-slave2 bash -c "mysql -u root -p1234 -e \"show slave status \G;\" " | grep Slave_IO_Running: | awk '{print $2}')
  slave2_SQL=$(docker exec -it mysql-slave2 bash -c "mysql -u root -p1234 -e \"show slave status \G;\" " | grep Slave_SQL_Running: | awk '{print $2}')
  echo "Slave_IO_Running: ${#slave2_IO} $slave2_IO"
  echo "Slave_SQL_Running: ${#slave2_SQL} $slave2_SQL"
  if ((${#slave2_IO} < 3 && ${#slave2_SQL} < 3)); then
    slave2copy
  fi
fi
