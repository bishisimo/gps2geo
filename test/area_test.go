//@Time : 2020/9/24 下午2:51
//@Author : bishisimo
package test

import (
	"github.com/bishisimo/errlog"
	geo "github.com/kellydunn/golang-geo"
	"gps2geo/geoBuilder"
	"gps2geo/utils"
	"os"
	"testing"
)

func BenchmarkWhereGps(b *testing.B) {
	areas := geoBuilder.GetAreas()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		areas.WhereGps(39.917544, 116.418757)
	}
}

func TestWhereGps(t *testing.T) {
	err := os.Setenv("RES_DIR", "../res")
	if errlog.Debug(err) {
		return
	}
	area := geoBuilder.GetAreas()
	gps := [][]float64{
		{24.643597, 117.943691},
		{21.536228, 107.972822},
		{30.908248, 120.437006},
		{37.999940, 100.918840},
	}
	if area.WhereGps(gps[3][0], gps[3][1]) == 0 {
		t.Fail()
	}
}
func TestCon(t *testing.T) {
	err := os.Setenv("RES_DIR", "../res")
	if errlog.Debug(err) {
		return
	}
	area := geoBuilder.GetAreas()
	gps := []*geo.Point{
		geo.NewPoint(30.908248, 120.437006),
		geo.NewPoint(37.999940, 100.918840),
	}
	utils.Println("0", area.Provinces[620000].ContainsPoint(gps[0]))
	utils.Println("1", area.Provinces[620000].Cities[620700].ContainsPoint(gps[0]))
	utils.Println("2", area.Provinces[620000].Cities[620700].Districts[620722].ContainsPoint(gps[0]))
}
