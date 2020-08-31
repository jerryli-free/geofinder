package service

import (
	"flag"
	"gpsService/lib"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	OSM_EXT  = "gps"
	OSM_LIST = "list.log"
)

type OsmFrm struct {
	OsmFolderPath  string
	OSM_POINT_MAP  map[int][]*Point
	OSM_INFO_MAP   map[int]*OsmDetail
	OSM_CENTER_MAP map[int]*Point
}

type OsmDetail struct {
	Osmid        int `json:"-"`
	Province_en  string
	Province_zh  string
	City_en      string
	City_zh      string
	District_en  string
	District_zh  string
	Similar      int
	Center_point *Point `json:"-""`
}
type Point struct {
	Lat float64
	Lon float64
}

/**
载入osm文件到内存
*/
func (g *OsmFrm) ReLoad() {
	map_path := g.OsmFolderPath + "/" + OSM_LIST
	if !lib.FileIsExist(map_path) {
		flag.Usage()
		panic("地理位置字典文件不存在：" + map_path)
	}
	g.loadList(map_path)
	//files, _ := ioutil.ReadDir(g.OsmFolderPath)
	//for _, file := range files {
	//加载osmdetail文件
	//	if ret, _ := regexp.MatchString(`^\d{4,15}\.`+OSM_EXT+`$`, strings.ToLower(file.Name())); ret {
	//g.loadosmfile(file.Name())
	//加载名称文件
	//	} else if file.Name() == OSM_LIST {
	//		g.loadList(g.OsmFolderPath + OSM_LIST)
	//	}
	//}
}

func (g *OsmFrm) GetOsmById(osmid int) *OsmDetail {
	return g.OSM_INFO_MAP[osmid]
}

func (g *OsmFrm) loadList(map_path string) {
	content, err := ioutil.ReadFile(map_path)
	if err != nil {
		panic(err)
	}
	//osm对应的pointlist
	g.OSM_POINT_MAP = make(map[int][]*Point)
	//osm对应的名称,中心点数据
	g.OSM_INFO_MAP = make(map[int]*OsmDetail)
	tmp_map := make(map[int]bool)
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimRight(line, "\r")
		arr := strings.Split(line, ",")
		if len(arr) >= 6 {
			osmid := 0
			if len(arr) == 6 {
				osmid, _ = strconv.Atoi(arr[5])
			} else if len(arr) == 9 {
				osmid, _ = strconv.Atoi(arr[8])
				pa_osmid, _ := strconv.Atoi(arr[5])
				tmp_map[pa_osmid] = true
			} else {
				continue
			}
			if len(arr) == 6 {
				g.OSM_INFO_MAP[osmid] = &OsmDetail{osmid, arr[0], arr[1], arr[3], arr[4], "", "", 0, &Point{}}
			} else {
				g.OSM_INFO_MAP[osmid] = &OsmDetail{osmid, arr[0], arr[1], arr[3], arr[4], arr[6], arr[7], 0, &Point{}}
			}
			g.loadosmfile(osmid)
		}
	}
	for key, _ := range tmp_map {
		delete(g.OSM_INFO_MAP, key)
		delete(g.OSM_POINT_MAP, key)
	}

}

/**
读取文件组建PointList，并且计算中心ponit
*/
func (g *OsmFrm) loadosmfile(osmid int) {
	filepath := g.OsmFolderPath + strconv.Itoa(osmid)+"." + OSM_EXT
	if !lib.FileIsExist(filepath) {
		panic("地理位置文件不存在：" + filepath)
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	//组建PointList
	g.OSM_POINT_MAP[osmid] = covertPoint(string(content))
	//计算中心点
	//g.calcCenterPoint(osmid)
	g.calcCenterPoint(osmid)
}

/**
获取不规则多边形重心点
*/
func (g *OsmFrm) calcCenterPoint(osmid int) {
	if _, ok := g.OSM_POINT_MAP[osmid]; !ok {
		return
	}
	points := g.OSM_POINT_MAP[osmid]
	length := len(points)
	area := 0.0 //多边形面积
	Gx := 0.0
	Gy := 0.0 // 重心的x、y
	for i := 1; i <= length; i++ {
		iLat := points[i%length].Lat
		iLng := points[i%length].Lon
		if iLat == 0 || iLng == 0 {
			continue
		}
		nextLat := points[i-1].Lat
		nextLng := points[i-1].Lon
		temp := (iLat*nextLng - iLng*nextLat) / 2.0
		area += temp
		Gx += temp * (iLat + nextLat) / 3.0
		Gy += temp * (iLng + nextLng) / 3.0
	}
	Gx = Gx / area
	Gy = Gy / area
	point := &Point{Gx, Gy}
	if _, ok := g.OSM_INFO_MAP[osmid]; !ok {
		g.OSM_INFO_MAP[osmid] = &OsmDetail{}
	}
	g.OSM_INFO_MAP[osmid].Center_point = point
}

/**
获取不规则多边形中心点
 */
func (g *OsmFrm) calcCenterPoint2(osmid int) {
	if _, ok := g.OSM_POINT_MAP[osmid]; !ok {
		return
	}
	points := g.OSM_POINT_MAP[osmid]
	length := len(points)
	x := 0.0
	y := 0.0
	z := 0.0

	for i := 0; i < length; i++ {
		if points[i].Lat == 0 || points[i].Lon == 0 {
			continue
		}
		lat := points[i].Lat * math.Pi / 180
		lon := points[i].Lon * math.Pi / 180

		a := math.Cos(lat) * math.Cos(lon)
		b := math.Cos(lat) * math.Sin(lon)
		c := math.Sin(lat)

		x += a
		y += b
		z += c
	}
	x /= float64(length)
	y /= float64(length)
	z /= float64(length)

	lon := math.Atan2(y, x)
	hyp := math.Sqrt(x*x + y*y)
	lat := math.Atan2(z, hyp)

	//return lat * 180 / math.Pi, lon * 180 / math.Pi
	Gx := lat * 180 / math.Pi
	Gy := lon * 180 / math.Pi
	point := &Point{Gx, Gy}
	if _, ok := g.OSM_INFO_MAP[osmid]; !ok {
		g.OSM_INFO_MAP[osmid] = &OsmDetail{}
	}
	g.OSM_INFO_MAP[osmid].Center_point = point
}

/**
将文件内容转换为PointList
*/
func covertPoint(content string) []*Point {
	ret := make([]*Point, 0)
	spaceRe, _ := regexp.Compile(`\s+`)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		ss := spaceRe.Split(line, -1)
		var lat, lon float64
		if len(ss) <= 1 {
			lat = 0
			lon = 0
		} else {
			lat, _ = strconv.ParseFloat(ss[2], 64)
			lon, _ = strconv.ParseFloat(ss[1], 64)
		}
		p := &Point{lat, lon}
		ret = append(ret, p)
	}
	return ret
}