//@Time : 2020/9/24 下午2:51
//@Author : bishisimo
package test

import (
	"github.com/smartystreets/goconvey/convey"
	"s2geo/geoBuilder"
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

	// Only pass t into top-level Convey calls
	convey.Convey("Given some integer with a starting value", t, func() {
		areas := geoBuilder.GetAreas()
		convey.Convey("The value should be greater by one", func() {
			convey.So(areas.WhereGps(39.917544, 116.418757), convey.ShouldEqual, 2)
		})
	})
}
