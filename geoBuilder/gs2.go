////@Time : 2020/9/23 上午9:45
////@Author : bishisimo
package geoBuilder

//import (
//	"fmt"
//	"github.com/golang/geo/s2"
//)
//
//func S2GeoPolygon(pio *PIO) {
//	s2ps:=make([]s2.Point,0)
//	coordinates:=pio.Geometry.Coordinates[0][0]
//	fmt.Println(pio.Properties.Name)
//	for i:=0;i<len(coordinates);i++{
//		//ll:=s2.LatLng{
//		//	Lat: s1.Angle(coordinates[i][1])*s1.Degree,
//		//	Lng: s1.Angle(coordinates[i][0])*s1.Degree,
//		//}
//		//s2ps = append(s2ps, s2.PointFromLatLng(s2.LatLngFromDegrees(coordinates[i][1]*s1.Degree.Radians(),coordinates[i][0]*s1.Degree.Radians())))
//		s2ps = append(s2ps, s2.Point{})
//	}
//	loop:=s2.LoopFromPoints(s2ps)
//	//fmt.Printf("%+v", loop)
//	gps:=[2]float64{118.4847810837,32.9595504398}
//	//isIn:=loop.ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(gps[1]*s1.Degree.Radians(),gps[0]*s1.Degree.Radians())))
//	isIn:=loop.ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(gps[1],gps[0])))
//	//fmt.Printf("\n%+v\n",s2.PointFromLatLng(s2.LatLngFromDegrees(gps[1]*s1.Degree.Radians(),gps[0]*s1.Degree.Radians())))
//	fmt.Printf("\n%+v\n",s2.PointFromLatLng(s2.LatLngFromDegrees(gps[1],gps[0])))
//	fmt.Printf("%+v\n", loop)
//	fmt.Println(isIn)
//}
