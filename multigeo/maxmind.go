package multigeo

import (
	"net"

	MM "github.com/oschwald/maxminddb-golang"
)

var mm *MM.Reader

type MaxMind struct {
	City struct {
		GeonameID int `maxminddb:"geoname_id"`
		Names     struct {
			EN string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"city"`
	Country struct {
		GeonameID int    `maxminddb:"geoname_id"`
		ISOCode   string `maxminddb:"iso_code"`
		Names     struct {
			EN string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"country"`
	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`
	Continent struct {
		GeonameID int    `maxminddb:"geoname_id"`
		Code      string `maxminddb:"code"`
		Names     struct {
			EN string `maxminddb:"en"`
		} `maxminddb:"names"`
	} `maxminddb:"continent"`
	Location struct {
		AccuracyRadius int     `maxminddb:"accuracy_radius"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		TimeZone       string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
}

func NewMaxMind(dbPath string) error {
	var err error
	mm, err = MM.Open(dbPath)
	if err != nil {
		return err
	}
	return nil
}

func (m *MaxMind) ToGeo(ip net.IP) (GeoResponse, error) {
	//var r MaxMind
	gr := GeoResponse{Provider: providerMaxMind}
	err := mm.Lookup(ip, &m)
	if err != nil {
		return gr, err
	}
	gr.City = m.City.Names.EN
	gr.Country = m.Country.Names.EN
	gr.CountryCode = m.Country.ISOCode
	gr.Continent = m.Continent.Names.EN
	//gr.ContinentCode = m.Continent.Code
	gr.ZipCode = m.Postal.Code
	gr.TimeZone = m.Location.TimeZone
	gr.Latitude = m.Location.Latitude
	gr.Longitude = m.Location.Longitude
	return gr, nil
}
