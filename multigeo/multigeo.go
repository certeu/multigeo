package multigeo

import (
	"net"
)

const (
	providerMaxMind     = "MaxMind"
	providerIP2Location = "IP2Location"
)

type GeoResponse struct {
	Provider    string  `json:"provider" xml:"provider"`
	City        string  `json:"city" xml:"city"`
	Country     string  `json:"country" xml:"country"`
	CountryCode string  `json:"countrycode" xml:"countrycode"`
	Continent   string  `json:"continent" xml:"continent"`
	Latitude    float64 `json:"latitude" xml:"latitude"`
	Longitude   float64 `json:"longitude" xml:"longitude"`
	ISP         string  `json:"isp" xml:"isp"`
	ZipCode     string  `json:"zipcode" xml:"zipcode"`
	TimeZone    string  `json:"timezone" xml:"timezone"`
	// other
	// International Direct Dialing
	IDDCode  string `json:"iddcode" xml:"ideecode"`
	AreaCode string `json:"arecode" xml:"areacode"`
	// WeatherStationCode
	WSCode string `json:"weatherstationcode" xml:"weatherstationcode"`
	// WeatherStationName
	WSName    string  `json:"weatherstationname" xml:"weatherstationname"`
	Elevation float64 `json:"elevation" xml:"elevation"`
	NetSpeed  string  `json:"netspeed" xml:"netspeed"`
	Domain    string  `json:"domain" xml:"domain"`
}

type Geoder interface {
	ToGeo(net.IP) (GeoResponse, error)
}
