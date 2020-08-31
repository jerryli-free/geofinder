# 项目介绍

使用go语言开发的通过经纬度取得城市省份的名称的微服务，用于当百度或者高德的免费服务出现请求限流的时候使用，在高并发或脚本批量运行使用时效果比较好



# 功能

v1.1版本目前可以查询国内的省份，直辖市，特区，和城市


# 安装

## 二进制安装

cd public;

go run build -o gps Main.go 



## Docker安装

docker pull registry.cn-hangzhou.aliyuncs.com/cb-repository/gpsservice:v1(私有仓库的密码咨询祖振)
docker run --name gpsservice -d  -p 1230:1230 -v /home/log/gpsServer:/home/log/gpsServer registry.cn-hangzhou.aliyuncs.com/cb-repository/gpsservice:v1



# 启动

执行二进制文件，nohup ./gps -osmpath ../service/osm/osmData/ &参数如下

  -acclog string
    	访问日志path（可不填，默认在/home/log/gpsServer/http_access.log）
  -osmpath string
    	osm文件目录 (default "osmData/"，必须项目，osm文件是程序加载的数据文件)
  -port int
    	指定端口号 (default 1230)

## 访问

curl '127.0.0.1:1230/?ac=GetGpsInfo&lat=39.884159&lon=117.010229'

## 返回值

{"errorCode":0,"errorMsg":"ok","result":{"Province_en":"Shandong","Province_zh":"山东省","City_en":"Dongying City","City_zh":"东营市","District_en":"Lijin County","District_zh":"利津县","Similar":1}}

errorCode:0代表成功，-1代表失败
errorMsg：错误内容，成功为ok

Province_en：省份的英文名称
Province_zh：省份的中文名称
City_en：城市英文
City_zh：城市中文
District_en：区县英文
District_zh：区县中文
Similar：是否为近似查找