# 项目介绍

使用go语言开发的通过经纬度取得城市省份的名称的微服务，20核的机器上qps可以达到10w，内存占用150M，在高并发或脚本批量运行使用时效果更加，在项目的public目前下，有一个已经编辑好的二进制文件，可以直接执行

```
./geofinder -osmpath ../service/osm/osmData/
```



# 功能

v2版本目前可以查询中国的省份，直辖市，特区，和城市


# 安装

## 二进制安装

```
cd public;

go build -o geofinder Main.go 
```

# 启动
cd public
go build -o geofinder Main.go 
执行二进制文件./geofinder 参数如下

```
-acclog string
    	访问日志path（可不填，默认在/home/log/geofinder/http_access.log）
-osmpath string
    	osm文件目录 (default "../service/osm/osmData/"，必须项目，osm文件是程序加载的数据文件)
-port int
    	指定端口号 (default 1230)
```



## 使用

curl '127.0.0.1:1230/getcity?lat=39.935297&lon=117.119176'

## 返回值说明

```json
{
    "errorCode": 0,//0代表成功，-1代表失败
    "errorMsg": "ok",//错误内容，成功为ok
    "result": {
        "Province_en": "Hebei",//省份的英文名称
        "Province_zh": "河北省",//省份的中文名称
        "City_en": "Langfang City",//城市英文
        "City_zh": "廊坊市",//城市中文
        "District_en": "Sanhe City",//区县英文
        "District_zh": "三河市",//区县中文
        "Similar": 0//1为准确查找，1为近似查找，就是查找距离输入点最近的城市
		"lat": "39.935297",
        "lon": "117.119176"
    }
}
```

