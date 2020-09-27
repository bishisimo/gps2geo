//@Time : 2020/9/23 上午9:45
//@Author : bishisimo
package geoBuilder

import (
	geo "github.com/kellydunn/golang-geo"
)

func GeoPolygon(pio PIO)*geo.Polygon {
	//info:=getInfo()
	ps:=make([]*geo.Point,0)
	coordinates:=pio.Geometry.Coordinates[0][0]
	for _,c:=range coordinates{
		ps = append(ps, geo.NewPoint(c[1],c[0]))
	}
	polygon:=geo.NewPolygon(ps)
	return polygon
}
func GeoPolygonNest(pio PIO)*[][]*geo.Polygon {
	ps:=make([]*geo.Point,0)
	polygons:=make([][]*geo.Polygon,0)
	coordinates:=pio.Geometry.Coordinates
	for i,c0:=range coordinates{
		polygons = append(polygons, make([]*geo.Polygon,0))
		for _,c1:=range c0{
			for _,c2:=range c1{
				ps = append(ps, geo.NewPoint(c2[1],c2[0]))
			}
			//polygon:=geo.NewPolygon(ps)
			polygons[i] = append(polygons[i],geo.NewPolygon(ps) )
		}
	}

	return &polygons
}