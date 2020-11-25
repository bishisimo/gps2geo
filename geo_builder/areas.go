//@Time : 2020/9/23 上午10:00
//@Author : bishisimo
package geo_builder

import (
	"fmt"
	"github.com/bishisimo/errlog"
	geo "github.com/kellydunn/golang-geo"
	"gps2geo/utils"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strings"
	"sync"
)

func NewAreas() *Areas {
	return &Areas{
		Provinces: make(map[int]*Province),
	}
}

type Areas struct {
	Provinces map[int]*Province
}

func (self *Areas) Add(code int, province *Province) {
	self.AddProvince(code, province)
}
func (self *Areas) AddProvince(code int, province *Province) {
	self.Provinces[code] = province
}
func (self *Areas) AddCity(code int, city *City) {
	codeProvince := code / 10000 * 10000
	province := self.Provinces[codeProvince]
	province.Add(code, city)
}
func (self *Areas) AddDistrict(code int, district *District) {
	codeProvince := code / 10000 * 10000
	province := self.Provinces[codeProvince]
	province.AddDistrict(code, district)
}

func (self Areas) WhereGps(lat float64, lng float64) int {
	pointP := geo.NewPoint(lat, lng)
	c := make(chan int, len(self.Provinces))
	wg := sync.WaitGroup{}
	var cityRecode *City
	//var provinceRecode *Province
	for _, province := range self.Provinces {
		wg.Add(1)
		go func(province *Province) {
			if province.ContainsPoint(pointP) {
				//provinceRecode = province
				for _, city := range province.Cities {
					wg.Add(1)
					go func(city *City) {
						if city.ContainsPoint(pointP) {
							cityRecode = city
							for _, district := range city.Districts {
								wg.Add(1)
								go func(district *District) {
									if district.ContainsPoint(pointP) {
										c <- district.AdCode
									}
									wg.Done()
								}(district)
							}
						}
						wg.Done()
					}(city)
				}
			}
			wg.Done()
		}(province)
	}
	wg.Wait()
	if len(c) > 0 {
		return <-c
	} else {
		adcode := self.WhereGpsInParticular(pointP)
		if adcode != 0 {
			return adcode
		} else if cityRecode != nil {
			adcode = cityRecode.WhereDistrictInApproximately(lat, lng)
			return adcode
		} else {
			return 0
		}
	}
}

func (self Areas) WhereGpsInParticular(pointP *geo.Point) int {
	c := make(chan int, len(self.Provinces))
	wg := sync.WaitGroup{}
	for _, province := range self.Provinces {
		wg.Add(1)
		go func(province *Province) {
			for _, city := range province.Cities {
				wg.Add(1)
				go func(city *City) {
					for _, district := range city.Districts {
						wg.Add(1)
						go func(district *District) {
							if district.ContainsPoint(pointP) {
								c <- district.AdCode
							}
							wg.Done()
						}(district)
					}
					//}
					wg.Done()
				}(city)
			}
			wg.Done()
		}(province)
	}
	wg.Wait()
	if len(c) == 1 {
		return <-c
	} else {
		return len(c)
	}
}

func NewProvince(properties *Properties, polygons *[][]*geo.Polygon) *Province {
	return &Province{
		Cities:       make(map[int]*City),
		MultiPolygon: polygons,
		Adcode:       properties.Adcode,
		Name:         properties.Name,
		ChildrenNum:  properties.ChildrenNum,
	}
}

type Province struct {
	Cities       map[int]*City
	MultiPolygon *[][]*geo.Polygon
	Adcode       int
	Name         string
	ChildrenNum  int
}

func (self *Province) Add(code int, city *City) {
	self.AddCity(code, city)
}
func (self *Province) AddCity(code int, city *City) {
	city.ProvinceName = self.Name
	if self.ChildrenNum == 0 {
		city.MultiPolygon = self.MultiPolygon
	}
	self.Cities[code] = city
}
func (self *Province) AddDistrict(code int, district *District) {
	codeCity := code / 100 * 100
	city := self.Cities[codeCity]
	if city == nil {
		city = self.Cities[code]
	}
	if city == nil {
		city = self.Cities[code/10000*10000]
	}
	city.Add(code, district)
}
func (self Province) ContainsPoint(pointP *geo.Point) bool {
	isIn := false
	for _, p0 := range *self.MultiPolygon {
		for i, p1 := range p0 {
			if p1.Contains(pointP) {
				isIn = i == 0
			}
		}
		if isIn {
			return isIn
		}
	}
	return false
}

func (self Province) ShowPolygon() {
	fmt.Printf("%v:\n", self.Name)
	fmt.Printf("shape:[%v,%v]:\n", len(*self.MultiPolygon), len((*self.MultiPolygon)[0]))
	fmt.Printf("data:%+v\n", *self.MultiPolygon)
}

func NewCity(properties *Properties, polygons *[][]*geo.Polygon) *City {
	return &City{
		Districts:    make(map[int]*District),
		MultiPolygon: polygons,
		ProvinceName: "",
		Adcode:       properties.Adcode,
		Name:         properties.Name,
		ChildrenNum:  properties.ChildrenNum,
	}
}

type City struct {
	Districts    map[int]*District
	MultiPolygon *[][]*geo.Polygon
	ProvinceName string
	Adcode       int
	Name         string
	ChildrenNum  int
}

func (self *City) Add(code int, district *District) {
	self.AddDistrict(code, district)
}
func (self *City) AddDistrict(code int, district *District) {
	district.ProvinceName = self.ProvinceName
	district.CityName = self.Name
	if self.ChildrenNum == 0 {
		district.MultiPolygon = self.MultiPolygon
	}
	self.Districts[code] = district
}
func (self City) ContainsPoint(pointP *geo.Point) bool {
	isIn := false
	for _, p0 := range *self.MultiPolygon {
		for i, p1 := range p0 {
			if p1.Contains(pointP) {
				isIn = i == 0
			}
		}
		if isIn {
			return isIn
		}
	}
	return false
}

func (self City) WhereDistrictInApproximately(lat float64, lng float64) int {
	for i := 3; i <= 10; i++ {
		offset := math.Pow(2, float64(i)) / math.Pow(10, 6)
		p := geo.NewPoint(lat-offset, lng)
		for _, districts := range self.Districts {
			if districts.ContainsPoint(p) {
				return districts.AdCode
			}
		}
		p = geo.NewPoint(lat+offset, lng)
		for _, districts := range self.Districts {
			if districts.ContainsPoint(p) {
				return districts.AdCode
			}
		}
		p = geo.NewPoint(lat, lng-offset)
		for _, districts := range self.Districts {
			if districts.ContainsPoint(p) {
				return districts.AdCode
			}
		}
		p = geo.NewPoint(lat, lng+offset)
		for _, districts := range self.Districts {
			if districts.ContainsPoint(p) {
				return districts.AdCode
			}
		}
	}
	return 0
}

func NewDistrict(properties *Properties, polygons *[][]*geo.Polygon) *District {
	return &District{
		Name:         properties.Name,
		MultiPolygon: polygons,
		ProvinceName: "",
		AdCode:       properties.Adcode,
		CityName:     "",
	}
}

type District struct {
	MultiPolygon *[][]*geo.Polygon
	ProvinceName string
	CityName     string
	AdCode       int
	Name         string
}

func (self District) ContainsPoint(pointP *geo.Point) bool {
	isIn := false
	for _, p0 := range *self.MultiPolygon {
		for i, p1 := range p0 {
			if p1.Contains(pointP) {
				isIn = i == 0
			}
		}
		if isIn {
			return isIn
		}
	}
	return false
}

var areas *Areas
var once sync.Once

func GetAreas() *Areas {
	once.Do(func() {
		areas = NewAreas()
		dirPath := os.Getenv("RES_DIR")
		if dirPath == "" {
			dirPath = "res"
		}
		files, err := ioutil.ReadDir(dirPath)
		if errlog.Debug(err) {
			return
		}
		var nationPath string
		provincesPaths := make(map[string]string, 0)
		citiesPath := make(map[string]string, 0)
		//将文件路径分别划分到省市
		for _, file := range files {
			name := strings.Split(file.Name(), ".")[0]
			code := strings.Split(name, "_")[0]
			filePath := path.Join(dirPath, file.Name())
			codeLen := len(code)
			municipality := map[string]string{"110000": "北京", "120000": "天津", "310000": "上海", "500000": "重庆"}
			if code == "100000" {
				nationPath = filePath
			} else if _, ok := municipality[code]; !ok && code[codeLen-4:] == "0000" {
				provincesPaths[code] = filePath
			} else {
				citiesPath[code] = filePath
			}
		}
		utils.Println("GeoJson init success")
		//构造省
		features := NewInfoFromJsonFile(nationPath).Features

		for _, pio := range features {
			province := NewProvince(&pio.Properties, GeoPolygonNest(&pio))
			areas.AddProvince(pio.Properties.Adcode, province)
			if strings.Contains(pio.Properties.Name, "市") {
				city := NewCity(&pio.Properties, GeoPolygonNest(&pio))
				areas.AddCity(pio.Properties.Adcode, city)
			}
		}
		//构造市
		for _, provincePath := range provincesPaths {
			features := NewInfoFromJsonFile(provincePath).Features
			for _, pio := range features {
				city := NewCity(&pio.Properties, GeoPolygonNest(&pio))
				areas.AddCity(pio.Properties.Adcode, city)
			}
		}
		//构造区
		for _, cityPath := range citiesPath {
			features := NewInfoFromJsonFile(cityPath).Features
			for _, pio := range features {
				district := NewDistrict(&pio.Properties, GeoPolygonNest(&pio))
				areas.AddDistrict(pio.Properties.Adcode, district)
			}
		}
	})
	return areas
}
