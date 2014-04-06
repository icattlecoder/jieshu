私人借书网站源码
=====================

online: [私人借书网](http://4jieshu.com)

私人借书网由开源go web开发框架[tgw](http://github.com/icattlecoder/tgw)开发，数据库采用MongoDB，数据从豆瓣阅读采集，同时采用豆瓣帐号登录。

本地测试运行：

```
~ go get github.com/icattlecoder/jieshu
~ cd src/github.com/icattlecoder/jieshu/www
~ go build
~ mongod 		#开启mongo服务
~ memcached -d 	#开启memcached服务
~ ./www 
```
nginx配置：

```
server {
  listen       80;
  server_name  www.jieshu.com;
  location / {
    proxy_pass http://127.0.0.1:8080;
  }
}
```
