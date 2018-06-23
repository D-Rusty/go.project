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

1. 添加git评论 仓库https://extremegtr.github.io/2017/09/07/Add-Gitment-comment-system-to-hexo-theme-NexT/
2. 头像上传是采用七牛存储还是腾讯云存储，再看，先使用默认图片
3. 错误提示采用页面标签的形式进行
4. 主页显示数据库中所有文章，在html中进行判断，用户在登录且文章作者和当前登录用户一致时显示 编辑，删除按钮
5. 准备部署外网