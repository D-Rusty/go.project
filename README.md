## 运行simpleBolog前请做如下配置：
1. 检查电脑环境的go_path，go_root路径
2. 执行以下命令，导入mysql包
```
go get github.com/go-sql-driver/mysql
go get github.com/mattn/go-sqlite3
go get github.com/mattn/go-sqlite3

```
3.确认程序没有错误时，执行bee run


# todo

3. 熟练理解beego 中mysql mongdb等操作
4. 了解 redis与beego关心

## simpleBlog部署腾讯云步骤
1. 申请腾讯云域名，并备案
2. 购买云主机并按照ubuntu系统
3.登入云主机,输入一下命令
```
 sudo apt-get update 更新系统自带的apt-get 
 wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz 从golang官网下载安装包
 sudo tar zxvf go1.10.3.linux-amd64.tar.gz -C /usr/local/ 解压golang包到/usr/local/目录
 mkdir /home/当前登录用户名/go.project 新建一个用于存放go项目源代码文件夹
 mkdir /home/当前登录用户名/go.project/src 新建一个所有go项目workspace
 vim ~/.bashrc 配置go环境变量，输入一下内容,输入完成后，按esc 输入！wq退出vim：
      export GOROOT=/usr/local/go
      export GOPATH=/home/当前登录用户名/go.project
      export GOBIN=/home/当前登录用户名/go.project/bin
      export PATH=$PATH:$GOPATH:/usr/local/go/bin:$GOBIN

 source ~/.bashrc 刷新 .bashrc
 go version 验证go是否配置成功，
 go env 验证go是否配置成功
 cd /home/当前登录用户名/go.project/src
 go get -u github.com/astaxie/beego 导入beego
 go get -u github.com/beego/bee 导入bee命令
 bee new hello 新建一个beego 工程
 sudo apt-get install ufw 按照端口管理工具
 sudo ufw allow 8080 开放8080端口给外网访问
 sudo ufw allow 80 开放80端口给外网访问
 sudo ufw status 检查端口是否已经开放
 cd hello  进入刚新建的beego 项目中
 bee run 运行刚新建的hello项目，通过云服务器公网IP：8080端口访问beego,如果能正常打开，则证明go环境，beego环境配置成功
 安装mysql，
  sudo apt-get install mysql-server
  sudo apt install mysql-client
  sudo apt install libmysqlclient-dev
 将数据库字符编码集修改为utf-8 https://www.jianshu.com/p/3111290b87f4
 mysql -uroot -p 登录 mysql
 create databases jikexueyuan 创建名字为极客学院的数据库
 sudo service mysql restart 重新启动数据库
 从github下载simeBlog源代码到/home/当前登录用户名/go.project/src下
 cd /home/当前登录用户名/go.project/src/simpleBlog
 按照以下simpleBlog运行依赖
 go get github.com/go-sql-driver/mysql
 go get github.com/mattn/go-sqlite3
 go get github.com/qiniu/api.v7/auth/qbox
 go get github.com/russross/blackfriday
  针对golang/net go get过程超时问题
 mkdir golang.org
 cd golang.org
 mkdir x
 cd x
 go clone https://github.com/golang/net.git
 go net install
 
 配置nginx
 sudo apt-get install nginx 按照nginx
 mkdir /home/当前登录用户名 xxxx.log
 sudo touch /etc/nginx/conf.d/beepkg.conf
 cd /etc/nginx/conf.d
 mkdir /home/当前登录用户名 xxxx.log
 sudo vim beepkg.conf并 beego.me中关于nginx发布进行配置按照，配置完成后退出
 cd /home/当前登录用户名/go.project/src/simpleBlog
 netstat -a 查询当前端口占用
 kill -9 pid 结束占用端口进程
 nohup bee run 运行simpleBlog
 sudo /etc/init.d/nginx reload 启动 nginx
```
将域名解析到云主机上 https://blog.csdn.net/LY_Dengle/article/details/78106594
可以开始域名访问simpleBlog了



