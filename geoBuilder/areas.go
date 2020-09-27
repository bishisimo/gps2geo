//@Time : 2020/9/23 上午10:00
//@Author : bishisimo
package geoBuilder

import (
	"fmt"
	"github.com/bishisimo/errlog"
	geo "github.com/kellydunn/golang-geo"
	"io/ioutil"
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
	province.Add(code/100*100, city)
}
func (self *Areas) AddDistrict(code int, district *District) {
	codeProvince := code / 10000 * 10000
	province := self.Provinces[codeProvince]
	province.AddDistrict(code, district)
}

func (self Areas) WhereGps(lat float64, lng float64) int {
	point := geo.NewPoint(lat, lng)
	c := make(chan int, len(self.Provinces))
	wg := sync.WaitGroup{}
	for _, province := range self.Provinces {
		wg.Add(1)
		go func(province *Province) {
			if province.ContainsPoint(point) {
				for _, city := range province.Cities {
					wg.Add(1)
					go func(city *City) {
						if city.ContainsPoint(point) {
							for _, district := range city.Districts {
								wg.Add(1)
								go func(district *District) {
									if district.ContainsPoint(point) {
										c <- district.Adcode
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
	if len(c) == 1 {
		return <-c
	} else if len(c) > 1 {
		return len(c)
	} else {
		return 0
	}
}

func NewProvince(properties *Properties, polygons *[][]*geo.Polygon) *Province {
	return &Province{
		Cities:      make(map[int]*City),
		Polygon:     polygons,
		Adcode:      properties.Adcode,
		Name:        properties.Name,
		ChildrenNum: properties.ChildrenNum,
	}
}

type Province struct {
	Cities      map[int]*City
	Polygon     *[][]*geo.Polygon
	Adcode      int
	Name        string
	ChildrenNum int
}

func (self *Province) Add(code int, city *City) {
	self.AddCity(code, city)
}
func (self *Province) AddCity(code int, city *City) {
	city.ProvinceName = self.Name
	if self.ChildrenNum == 0 {
		city.Polygon = self.Polygon
	}
	self.Cities[code] = city
}
func (self *Province) AddDistrict(code int, district *District) {
	codeCity := code / 100 * 100
	city := self.Cities[codeCity]
	//if city == nil {
	//	fmt.Println(code)
	//	city = self.Cities[code]
	//}
	if city == nil {
		city = self.Cities[code/10000*10000]
	}
	city.Add(code, district)
}
func (self Province) ContainsPoint(pointP *geo.Point) bool {
	for _, polygon := range *self.Polygon {
		for _, p1 := range polygon {
			if p1.Contains(pointP) {
				return true
			}
		}
	}
	return false
}

func (self Province) ShowPolygon() {
	fmt.Printf("%v:\n", self.Name)
	fmt.Printf("shape:[%v,%v]:\n", len(*self.Polygon), len((*self.Polygon)[0]))
	fmt.Printf("data:%+v\n", *self.Polygon)
}

func NewCity(properties *Properties, polygons *[][]*geo.Polygon) *City {
	return &City{
		Districts:    make(map[int]*District),
		Polygon:      polygons,
		ProvinceName: "",
		Adcode:       properties.Adcode,
		Name:         properties.Name,
		ChildrenNum:  properties.ChildrenNum,
	}
}

type City struct {
	Districts    map[int]*District
	Polygon      *[][]*geo.Polygon
	ProvinceName string
	Adcode       int
	Name         string
	ChildrenNum  int
}

func (self *City) Add(code int, district *District) {
	self.AddDistrict(code, district)
}
func (self *City) AddDistrict(code int, district *District) {
	//defer func() {
	//	if err:=recover();err!=nil{
	//		fmt.Printf("AddDistrict:%+v\n",self)
	//	}
	//}()
	district.ProvinceName = self.ProvinceName
	district.CityName = self.Name
	if self.ChildrenNum == 0 {
		district.Polygon = self.Polygon
	}
	self.Districts[code] = district
}
func (self City) ContainsPoint(pointP *geo.Point) bool {
	for _, polygon := range *self.Polygon {
		for _, p1 := range polygon {
			if p1.Contains(pointP) {
				return true
			}
		}
	}
	return false
}
func NewDistrict(properties *Properties, polygons *[][]*geo.Polygon) *District {
	return &District{
		Name:         properties.Name,
		Polygon:      polygons,
		ProvinceName: "",
		Adcode:       properties.Adcode,
		CityName:     "",
	}
}

type District struct {
	Polygon      *[][]*geo.Polygon
	ProvinceName string
	CityName     string
	Adcode       int
	Name         string
}

func (self District) ContainsPoint(pointP *geo.Point) bool {
	for _, polygon := range *self.Polygon {
		for _, p1 := range polygon {
			if p1.Contains(pointP) {
				return true
			}
		}
	}
	return false
}

var areas *Areas
var once sync.Once

func GetAreas() *Areas {
	//var once sync.Once

	once.Do(func() {
		areas = NewAreas()
		dirPath := "res"
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
		fmt.Printf("%v\n", nationPath)
		//构造省
		features := NewInfoFromJsonFile(nationPath).Features

		for _, pio := range features {
			if pio.Properties.Adcode == 0 {
				fmt.Println("Adcode", nationPath, pio.Properties.Name)
			}
			province := NewProvince(&pio.Properties, GeoPolygonNest(pio))
			areas.AddProvince(pio.Properties.Adcode, province)
			if strings.Contains(pio.Properties.Name, "市") {
				city := NewCity(&pio.Properties, GeoPolygonNest(pio))
				areas.AddCity(pio.Properties.Adcode, city)
			}
		}
		//构造市
		for _, provincePath := range provincesPaths {
			features := NewInfoFromJsonFile(provincePath).Features
			for _, pio := range features {
				city := NewCity(&pio.Properties, GeoPolygonNest(pio))
				areas.AddCity(pio.Properties.Adcode, city)
			}
		}
		//fmt.Printf("%+v", areas.Provinces)

		//构造区
		for _, cityPath := range citiesPath {
			features := NewInfoFromJsonFile(cityPath).Features
			for _, pio := range features {
				if pio.Properties.Adcode == 0 {
					fmt.Println("Adcode", cityPath, pio.Properties.Name)
				}
				district := NewDistrict(&pio.Properties, GeoPolygonNest(pio))
				areas.AddDistrict(pio.Properties.Adcode, district)
			}
		}
	})
	return areas
}
