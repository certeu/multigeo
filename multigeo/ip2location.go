package multigeo

import (
	"fmt"
	"net"
	"os"

	IP2L "github.com/ip2location/ip2location-go"
)

type IP2Location struct {
	ip2l *IP2L.IP2Locationrecord
}

func NewIP2Location(dbPath string) error {
	_, err := os.Open(dbPath)
	if err != nil {
		return err
	}
	IP2L.Open(dbPath)
	return nil
}

func (i IP2Location) ToGeo(ip net.IP) (GeoResponse, error) {
	res := IP2L.Get_all(fmt.Sprintf("%s", ip))
	gr := GeoResponse{
		Provider:    providerIP2Location,
		City:        res.City,
		Country:     res.Country_long,
		CountryCode: res.Country_short,
		Latitude:    float64(res.Latitude),
		Longitude:   float64(res.Longitude),
		ISP:         res.Isp,
		ZipCode:     res.Zipcode,
		TimeZone:    res.Timezone,
		IDDCode:     res.Iddcode,
		AreaCode:    res.Areacode,
		WSCode:      res.Weatherstationcode,
		WSName:      res.Weatherstationname,
		Elevation:   float64(res.Elevation),
		NetSpeed:    res.Netspeed,
		Domain:      res.Domain,
	}
	return gr, nil
}

func init() {
	// IP2Location doesn't return anything when opening a database.
	// We'll load the database and use package functions
}
