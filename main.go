//@Time : 2020/9/22 上午11:03
//@Author : bishisimo
package main

import (
	"github.com/gin-gonic/gin"
	"gps2geo/geoBuilder"
	"strconv"
)

func web(areas *geoBuilder.Areas) {
	r := gin.Default()
	r.GET("/gps/where", func(c *gin.Context) {
		lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
		lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
		c.JSON(200, gin.H{
			"data": areas.WhereGps(lat,lng),
		})
	})

	_ = r.Run(":8100")
}

func main() {
	areas := geoBuilder.GetAreas()
	web(areas)
}
