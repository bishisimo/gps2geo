//@Time : 2020/9/22 上午11:03
//@Author : bishisimo
package main

import (
	"github.com/gin-gonic/gin"
	"gps2geo/geo_builder"
	"strconv"
)

func web(areas *geo_builder.Areas) {
	r := gin.Default()
	r.GET("/gps/where", func(c *gin.Context) {
		lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
		lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
		c.JSON(200, gin.H{
			"data": areas.WhereGps(lat, lng),
		})
	})

	_ = r.Run(":8100")
}

func main() {
	areas := geo_builder.GetAreas()
	web(areas)
}
