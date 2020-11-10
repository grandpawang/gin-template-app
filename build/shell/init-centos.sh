# 设置主机名
hostnamectl set-hostname coint.server

mkdir /.old/

# 设置合盖不休眠
sed -i "s/#HandleLidSwitch=suspend/HandleLidSwitch=ignore/g" /etc/systemd/logind.conf
systemctl restart systemd-logind

# 关闭防火墙
setenforce 0
sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config

systemctl stop firewalld
systemctl disable firewalld
systemctl stop iptables
systemctl disable iptables

systemctl status firewalld
systemctl status iptables

# 配置yum源
#curl 下载阿里镜像源
cp /etc/yum.repos.d/CentOS-Base.repo /.old/
cp /etc/yum.repos.d/epel-7.repo /.old/
wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
wget -O /etc/yum.repos.d/epel-7.repo http://mirrors.aliyun.com/repo/epel-7.repo
yum clean all && yum makecache
yum -y update && yum -y upgrade

# 安装基础工具
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
yum install -y net-tools \
    vim \
    wget \
    lsof \
    yum-utils \
    device-mapper-persistent-data \
    lvm2 \
    python-pip \
    docker-ce \
    docker-compose

# 关闭swap
swapoff -a
row=$(nl /etc/fstab | grep ext4 | awk '{print $1}')
if [ -n "$row" ]; then
    echo "need to delete swap"
    echo "sed -i "$row"'d' /etc/fstab"
    cat /etc/fstab
    cp /etc/fstab /.old/fstab.tmp
    sed -i "$row"'d' /etc/fstab
fi
free -m

# 安装docker
sudo cat <<EOF >/etc/docker/deamon.json
{
	"registry-mirrors": [
			"https://almtd3fa.mirror.aliyuncs.com",
			"http://registry.docker-cn.com",
			"http://docker.mirrors.ustc.edu.cn",
			"http://hub-mirror.c.163.com"
	],
	"insecure-registries": [],
	"debug": true,
	"experimental": true
}
EOF
sudo usermod -aG docker root
systemctl daemon-reload
systemctl restart docker.service
sudo systemctl enable docker
sudo systemctl start docker

# 安装golang
wget https://studygolang.com/dl/golang/go1.14.4.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.14.4.linux-amd64.tar.gz
cat >>/etc/profile <<EOF
export GOROOT=/usr/local/go
export PATH=\$PATH:\$GOROOT/bin
export GOPATH=/root/go
export PATH=\$PATH:\$GOPATH/BIN
EOF
source /etc/profile
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.io,direct

# 安装nodejs
wget https://nodejs.org/dist/v12.18.2/node-v12.18.2-linux-x64.tar.xz
tar -C /usr/local -xvf node-v12.18.2-linux-x64.tar.xz
mw /usr/local/node-v12.18.2-linux-x64/ /usr/local/node/
cat >>/etc/profile <<EOF
export PATH=\$PATH:/usr/local/node/bin
EOF
source /etc/profile
npm config set registry https://registry.npm.taobao.org

# 安装宝塔
wget -O install.sh http://download.bt.cn/install/install_6.0.sh && sh install.sh
