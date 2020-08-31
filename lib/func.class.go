package lib

import (
	"math"
	"os"
)

const (
	LAT_MIN = -85.05112878
	LAT_MAX = 85.05112878
	LON_MIN = -180
	LON_MAX = 180
)

func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
func GetDistance(latitude1, longitude1, latitude2, longitude2 float64) float64 {
	//计算经纬度，返回单位是 米
	earth_radius := 6378.137
	radLat1 := radius(latitude1)
	radLat2 := radius(latitude2)
	a := radLat1 - radLat2
	b := radius(longitude1) - radius(longitude2)
	s := 2 * math.Asin((math.Sqrt(math.Pow(math.Sin(a/2), 2) + math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2))))
	s = s * earth_radius
	s = math.Round(s*10000) / 10
	return s
}
func radius(d float64) float64 {
	return d * math.Pi / 180.0
}
func Interleave64(xlo float64, ylo float64) int64 {
	B := []int{0x5555555555555555, 0x3333333333333333,
		0x0F0F0F0F0F0F0F0F, 0x00FF00FF00FF00FF,
		0x0000FFFF0000FFFF}
	S := []int{1, 2, 4, 8, 16}

	x := int(xlo)
	y := int(ylo)

	x = (x | (x << S[4])) & B[4]
	y = (y | (y << S[4])) & B[4]

	x = (x | (x << S[3])) & B[3]
	y = (y | (y << S[3])) & B[3]

	x = (x | (x << S[2])) & B[2]
	y = (y | (y << S[2])) & B[2]

	x = (x | (x << S[1])) & B[1]
	y = (y | (y << S[1])) & B[1]

	x = (x | (x << S[0])) & B[0]
	y = (y | (y << S[0])) & B[0]
	return int64(x | (y << 1))
}

func Deinterleave64(interleaved int64) int64 {
	B := []int64{0x5555555555555555, 0x3333333333333333,
		0x0F0F0F0F0F0F0F0F, 0x00FF00FF00FF00FF,
		0x0000FFFF0000FFFF, 0x00000000FFFFFFFF}
	S := []int{0, 1, 2, 4, 8, 16}

	x := interleaved
	y := interleaved >> 1

	x = (x | (x >> S[0])) & B[0]
	y = (y | (y >> S[0])) & B[0]

	x = (x | (x >> S[1])) & B[1]
	y = (y | (y >> S[1])) & B[1]

	x = (x | (x >> S[2])) & B[2]
	y = (y | (y >> S[2])) & B[2]

	x = (x | (x >> S[3])) & B[3]
	y = (y | (y >> S[3])) & B[3]

	x = (x | (x >> S[4])) & B[4]
	y = (y | (y >> S[4])) & B[4]

	x = (x | (x >> S[5])) & B[5]
	y = (y | (y >> S[5])) & B[5]

	return x | (y << 32)
}

func GeohashDecode(hash int64, step int) (float64, float64) {
	hashDeinterleave := Deinterleave64(hash)
	latBits := hashDeinterleave & 4294967295 //低32位
	lonBits := hashDeinterleave >> 32        //高32位
	x := int64(1 << step)
	longitude := float64(LON_MIN) + float64(float64(lonBits)/float64(x)*float64(LON_MAX-LON_MIN)) //经度
	latitude := float64(LAT_MIN) + float64(latBits)/float64(x)*float64(LAT_MAX-LAT_MIN)           //纬度
	return longitude, latitude
}

func GeohashEncode(lon float64, lat float64, step int) int64 {
	latOffset := (lat - LAT_MIN) / (LAT_MAX - LAT_MIN)
	lonOffset := (lon - LON_MIN) / (LON_MAX - LON_MIN)
	x := 1 << step
	latOffset *= float64(x)
	lonOffset *= float64(x)
	return Interleave64(latOffset, lonOffset)
}
