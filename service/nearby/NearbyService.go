package service

import (
	"gpsService/lib"
	osmService "gpsService/service/osm"
	"sort"
)

const NEARBY_BIT = 26

type NearBy struct {
	SortList []Nearbydetail
}
type Nearbydetail struct {
	Osmid    int
	hashcode int64
}
type nearbydetaillist []Nearbydetail

func (s nearbydetaillist) Len() int           { return len(s) }
func (s nearbydetaillist) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s nearbydetaillist) Less(i, j int) bool { return s[i].hashcode < s[j].hashcode }

func (n *NearBy) Reload(osmdetails map[int]*osmService.OsmDetail) {
	for _, osmdetail := range osmdetails {
		hashcode := lib.GeohashEncode(osmdetail.Center_point.Lon, osmdetail.Center_point.Lat, NEARBY_BIT)
		n.SortList = append(n.SortList, Nearbydetail{Osmid: osmdetail.Osmid, hashcode: hashcode})
	}
	sort.Sort(nearbydetaillist(n.SortList))
}

func (n *NearBy) GetNearbyCity(lat float64, lon float64, radius int, middleIndex int) ([]Nearbydetail,int) {
	hashcode := lib.GeohashEncode(lon, lat, NEARBY_BIT)
	ret,middleIndex := n.BinaryFind(hashcode, radius, middleIndex)
	return ret,middleIndex
}

func (n *NearBy) BinaryFind(findVal int64, radius int, middleIndex int) ([]Nearbydetail,int) {

	leftIndex := 0
	rightIndex := len(n.SortList) - 1
	if middleIndex == -1 {
		for {
			if leftIndex > rightIndex {
				break
			}
			middleIndex = (leftIndex + rightIndex) / 2
			if n.SortList[middleIndex].hashcode == findVal {
				break
			} else if n.SortList[middleIndex].hashcode > findVal {
				rightIndex = middleIndex - 1
			} else if n.SortList[middleIndex].hashcode < findVal {
				leftIndex = middleIndex + 1
			}
		}
	}
	middle_start := 0
	if middleIndex-radius > middle_start {
		middle_start = middleIndex - radius
	}
	middle_end := len(n.SortList)
	if middleIndex+radius < middle_end {
		middle_end = middleIndex + radius
	}
	//fmt.Println(n.SortList,middle_start,middle_end,)
	return n.SortList[middle_start:middle_end],middleIndex

}
