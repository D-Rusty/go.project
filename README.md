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
2. 头像上传是采用七牛存储还是腾讯云存储，再看，先使用默认图片
3. 错误提示采用页面标签的形式进行
4. 准备部署外网
5. 添加markdown格式兼容