package app

import (
	"encoding/json"
	"github.com/fvbock/endless"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"gpsService/lib"
	service "gpsService/service/gpsfind"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

const (
	RET_OK  = 0
	RET_ERR = 1
)

type GpsServer struct {
	Port          int
	AccessLogPath string
	ErrLogPath    string
	GpsService    *service.GpsService
}

func (g *GpsServer) Start() {
	app := gin.New()
	app.Use(lib.NewUseLog(g.AccessLogPath), gin.Recovery())
	app.Use(gzip.Gzip(gzip.DefaultCompression)) //开启gzip压缩. 请求头中含 "Accept-Encoding : gzip" 才会压缩
	g.GpsService.Init()
	app.GET("/", g.routeGpsServer)
	endless.ListenAndServe(":"+strconv.Itoa(g.Port), app)
}

func (g *GpsServer) routeGpsServer(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			lib.Fatal(r)
			echoResult(RET_ERR, "系统异常", nil, ctx)
		}
	}()

	beginTime := time.Now() //计时
	var ret interface{}
	errorMsg := "ok"
	errorCode := RET_OK
	lat_str := ctx.Request.FormValue("lat")
	lon_str := ctx.Request.FormValue("lon")
	lat, _ := strconv.ParseFloat(lat_str, 64)
	lon, _ := strconv.ParseFloat(lon_str, 64)
	methodName := ctx.Request.FormValue("ac")

	switch methodName {
	case "GetGpsInfo":
		ret = g.GpsService.FindByGps(lat, lon)
		vi := reflect.ValueOf(ret)
		if vi.IsNil() {
			lib.Info("notfind lon:" + strconv.FormatFloat(lon, 'E', -1, 64) + " lat:" + strconv.FormatFloat(lat, 'E', -1, 64))
		}
		break
	default:
		errorCode = RET_ERR
		errorMsg = "接口不存在"
	}

	acElapsed := time.Since(beginTime) //接口耗时
	echoResult(errorCode, errorMsg, ret, ctx)
	totalElapsed := time.Since(beginTime) //总耗时
	if totalElapsed > 500*time.Millisecond {
		//超过500毫秒
		retJson, _ := json.Marshal(ret)
		logMsg := "Elapsed Time: lat:" + lat_str + "|lon:" + lon_str + "|acElapsed:" + acElapsed.String() + "|totalElapsed:" + totalElapsed.String()
		logMsg += "|url:" + ctx.Request.RequestURI
		logMsg += "|ret:" + string(retJson)
		lib.SetInfoPath("/home/log/gpsServer/slow_access.log").Info("notfind lon:" + strconv.FormatFloat(lon, 'E', -1, 64) + " lat:" + strconv.FormatFloat(lat, 'E', -1, 64))
	}
}

func echoResult(errcode int, errmsg string, result interface{}, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"errorCode": errcode,
		"errorMsg":  errmsg,
		"result":    result,
	})
}
