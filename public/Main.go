package main

import (
	"flag"
	"gpsService/app"
	service "gpsService/service/gpsfind"
)

func main() {
	var port int
	var osmpath,acclog string
	flag.IntVar(&port, "port", 1230, "指定端口号")
	flag.StringVar(&osmpath, "osmpath", "../service/osm/osmData/", "osm文件目录")
	flag.StringVar(&acclog, "acclog", "", "访问日志path")
	flag.Parse()
	serv := &app.GpsServer{
		Port:          port,
		AccessLogPath: acclog,
		GpsService:    &service.GpsService{OsmFolder: osmpath},
	}
	serv.Start()

}

