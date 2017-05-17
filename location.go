package geoip

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
)

// Location represents the GeoLite City Location csv.
type Location struct {
	// LocationID is the locations id in the db.
	LocationID int `json:"locId"`
	// CountryCode is a two-character ISO 3166-1 country code for
	// the country associated with the IP address.
	CountryCode string `json:"country"`
	// RegionCode is a two character ISO-3166-2 or FIPS 10-4 code
	// for the state/region associated with the IP address.
	RegionCode string `json:"region"`
	// City is the city or town name associated with the IP address.
	City string `json:"city"`
	// PostalCode is the postal code associated with the IP address.
	PostalCode string `json:"postalCode"`
	// Latitude is the latitude associated with the IP address
	Latitude float32 `json:"latitude"`
	// Longitude is the longitude associated with the IP address.
	Longitude float32 `json:"longitude"`
	// MetroCode is the metro code associated with the IP address.
	MetroCode uint16 `json:"metroCode"`
	// AreaCode is the telephone area code associated with the IP address
	AreaCode string `json:"areaCode"`
	// Serialized is the location serialized into JSON
	Serialized []byte `json:"-"`
}

// LocationMap maps a location id with the location.
type LocationMap struct {
	data map[int]*Location
}

// NewLocation will parse a csv record into a Location struct. It assumes the following order
// locId,country,region,city,postalCode,latitude,longitude,metroCode,areaCode
func NewLocation(record []string) (*Location, error) {
	var err error
	loc := &Location{
		CountryCode: record[1],
		RegionCode:  record[2],
		City:        record[3],
		PostalCode:  record[4],
		AreaCode:    record[8],
	}
	loc.LocationID, err = strconv.Atoi(record[0])
	if err != nil {
		return nil, err
	}
	if record[5] != "" {
		lat, err := strconv.ParseFloat(record[5], 32)
		if err != nil {
			return nil, err
		}
		loc.Latitude = float32(lat)
	}
	if record[6] != "" {
		long, err := strconv.ParseFloat(record[6], 32)
		if err != nil {
			return nil, err
		}
		loc.Longitude = float32(long)
	}
	if record[7] != "" {
		metroCode, err := strconv.Atoi(record[7])
		if err != nil {
			return nil, err
		}
		loc.MetroCode = uint16(metroCode)
	}

	loc.Serialized, err = json.Marshal(loc)
	if err != nil {
		return nil, err
	}
	return loc, nil
}

// NewLocationMap will load a location csv from file and return a
// mapping of location id to location.
func NewLocationMap(f io.Reader) (*LocationMap, error) {
	lm := &LocationMap{
		data: map[int]*Location{},
	}

	r := csv.NewReader(f)
	r.FieldsPerRecord = 9
	fmt.Println("Locations:")
	i := -1
READ_LOOP:
	for {
		i++
		record, err := r.Read()
		if err == io.EOF {
			break READ_LOOP
		}
		if err != nil {
			if perr, ok := err.(*csv.ParseError); ok && perr.Err == csv.ErrFieldCount {
				if i != 0 {
					log.Println(err)
				}
				continue READ_LOOP
			}
			return nil, err
		}
		if len(record) == 9 {
			loc, err := NewLocation(record)
			if err != nil {
				log.Println(err)
				continue READ_LOOP
			}
			lm.data[loc.LocationID] = loc
			fmt.Printf("\r %d", i)
		}
	}
	fmt.Printf("\ndone\n\n")
	return lm, nil
}

// LookupByID will lookup a location by id. This assumes
func (l *LocationMap) LookupByID(id int) (*Location, bool) {
	loc, ok := l.data[id]
	return loc, ok
}
