package service

import (
	"gpsService/lib"
	nearbyService "gpsService/service/nearby"
	osmService "gpsService/service/osm"
	"math"
)

type GpsService struct {
	Osm       *osmService.OsmFrm
	Nearby    *nearbyService.NearBy
	OsmFolder string
}

func (g *GpsService) Init() {
	g.Osm = &osmService.OsmFrm{}
	g.Osm.OsmFolderPath = g.OsmFolder
	g.Osm.ReLoad()
	g.Nearby = &nearbyService.NearBy{}
	g.Nearby.Reload(g.Osm.OSM_INFO_MAP)
}
func (g *GpsService) isInCity(point *osmService.Point, osmid int) bool {
	polygonPoints := g.Osm.OSM_POINT_MAP[osmid]
	flag := false
	j := len(polygonPoints) - 1
	for i := 0; i < len(polygonPoints); i++ {
		p1 := polygonPoints[i]
		p2 := polygonPoints[j]
		j = i
		//剔除间隔点
		if p1.Lat == 0 || p2.Lat == 0 || p1.Lon == 0 || p2.Lon == 0 {
			continue
		}
		//剔除不在范围内的
		if point.Lon < math.Min(p1.Lon, p2.Lon) || point.Lon > math.Max(p1.Lon, p2.Lon) {
			continue
		}
		// 如果和某一个共点那么直接返回true
		if point.Lon == p1.Lon && point.Lat == p1.Lat {
			return true
		}
		if point.Lon == p2.Lon && point.Lat == p2.Lat {
			return true
		}
		// 如果和两点共线
		if point.Lon == p1.Lon && point.Lon == p2.Lon {
			if point.Lat < math.Min(p1.Lat, p2.Lat) || point.Lon > math.Max(p1.Lat, p2.Lat) {
				return false
			} else {
				return true
			}
		}
		// 这里判断是否刚好被测点在多边形的边上 TODO
		//如果穿过两交点 TODO
		//查看点是否在线的左侧
		if (p1.Lon > point.Lon) != (p2.Lon > point.Lon) && (point.Lat < (point.Lon-p1.Lon)*(p1.Lat-p2.Lat)/(p1.Lon-p2.Lon)+p1.Lat) {
			flag = !flag
		}
	}
	return flag
}

func (g *GpsService) FindByGps(lat float64, lon float64) *osmService.OsmDetail {
	var m *osmService.OsmDetail
	curpoint := &osmService.Point{lat, lon}
	radius_arr := []int{2, 5, 10, 20, 50}
	var nears []nearbyService.Nearbydetail
	middleIndex := -1
out:
	for _, radius := range radius_arr {
		searched_map := make(map[int]bool)
		nears, middleIndex = g.Nearby.GetNearbyCity(lat, lon, radius, middleIndex)
		//for _, near := range nears {
		//	fmt.Println(near.Osmid, g.Osm.OSM_INFO_MAP[near.Osmid])
		//}

		for _, near := range nears {
			if _, ok := searched_map[near.Osmid]; !ok {
				if g.isInCity(curpoint, near.Osmid) {
					m = g.Osm.OSM_INFO_MAP[near.Osmid]
					break out
				}
				searched_map[near.Osmid] = true
			}
		}
	}

	if m == nil {
		var min float64
		var osmid int
		for _, near := range nears {
			distance := lib.GetDistance(lat, lon, g.Osm.OSM_INFO_MAP[near.Osmid].Center_point.Lat, g.Osm.OSM_INFO_MAP[near.Osmid].Center_point.Lon)
			//fmt.Println(distance, near.Osmid)
			if min == 0 || min > distance {
				min = distance
				osmid = near.Osmid
			}
		}
		m = g.Osm.OSM_INFO_MAP[osmid]
		m.Similar = 1
	}

	return m
}
