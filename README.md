# qqwry
> IP 地址归属地查询库, 使用纯真IP库

## 安装
```shell
go get github.com/mazezen/qqwry
```

### 获取最新的纯真 IP 库
访问 http://www.cz88.net 下载纯真IP库.或者使用如下脚本
```shell
php download.php

./download.sh
```

### 使用方法
```go
qqWry, err := NewQQWry("./qqwry.dat")
if err != nil {
	panic(err)
}
ip, country, area := qqWry.Find("8.8.8.8")
fmt.Printf("ip: %s, country: %s, area: %s", ip, country, area)
```


### 特别感谢
感谢<a href="https://www.cz88.net/">纯真IP库</a>一直坚持为大家提供免费IP数据库。

感谢<a href="https://github.com/yinheli">yinheli</a>的<a href="https://github.com/yinheli/qqwry">qqwry</a>项目，为我提供纯真ip库解析算法参考。

感谢<a href="https://github.com/metowolf">metowolf</a>的<a href="https://github.com/metowolf/qqwry.ipdb">qqwry.ipdb</a>项目，提供纯真czdb转ipdb数据库。
