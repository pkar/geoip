package geoip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var (
	Fonda    *Location
	Laval    *Location
	Research *Location
)

func init() {
	Fonda = &Location{
		LocationID:  24107,
		CountryCode: "US",
		RegionCode:  "NY",
		City:        "Fonda",
		PostalCode:  "12068",
		Latitude:    float32(42.9508),
		Longitude:   float32(-74.3937),
		MetroCode:   532,
		AreaCode:    "518",
	}
	Fonda.Serialized, _ = json.Marshal(Fonda)

	Laval = &Location{
		LocationID:  24100,
		CountryCode: "CA",
		RegionCode:  "QC",
		City:        "Laval",
		PostalCode:  "h7w4s8",
		Latitude:    float32(45.6167),
		Longitude:   float32(-73.7500),
		MetroCode:   0,
		AreaCode:    "",
	}
	Laval.Serialized, _ = json.Marshal(Laval)

	Research = &Location{
		LocationID:  609013,
		CountryCode: "AU",
		RegionCode:  "07",
		City:        "Research",
		PostalCode:  "3095",
		Latitude:    float32(-37.7000),
		Longitude:   float32(145.1833),
		MetroCode:   0,
		AreaCode:    "",
	}
	Research.Serialized, _ = json.Marshal(Research)
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func Test_locationHandlerLocationIDLookup(t *testing.T) {
	r := strings.NewReader(`garbage
1,"O1","","","",0.0000,0.0000,,
24107,"US","NY","Fonda","12068",42.9508,-74.3937,532,518
`)
	lm, err := NewLocationMap(r)
	ok(t, err)
	locHandler := &locationHandler{
		LocationMap: lm,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "localhost:8899?locId=24107", nil)
	locHandler.ServeHTTP(w, req)
	equals(t, 200, w.Code)
	equals(t, string(Fonda.Serialized), w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "localhost:8899?locId=44", nil)
	locHandler.ServeHTTP(w, req)
	equals(t, 404, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "localhost:8899?locId=💩", nil)
	locHandler.ServeHTTP(w, req)
	equals(t, 400, w.Code)
}

func Test_locationHandlerIPLookup(t *testing.T) {
	r := strings.NewReader(`garbage
1,"O1","","","",0.0000,0.0000,,
24107,"US","NY","Fonda","12068",42.9508,-74.3937,532,518
609013,"AU","07","Research","3095",-37.7000,145.1833,,
`)
	lm, err := NewLocationMap(r)
	ok(t, err)

	rb := strings.NewReader(`garbage
"16777216","16777471","609013"
"16777472","16778239","104084"
`)
	bm, err := NewBlocksMap(rb)
	ok(t, err)
	locHandler := &locationHandler{
		LocationMap: lm,
		BlocksMap:   bm,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "localhost:8899?ip=1.0.0.255", nil)
	locHandler.ServeHTTP(w, req)
	equals(t, 200, w.Code)
	equals(t, string(Research.Serialized), w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "localhost:8899?ip=44", nil)
	locHandler.ServeHTTP(w, req)
	equals(t, 404, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "localhost:8899?ip=💩", nil)
	locHandler.ServeHTTP(w, req)
	equals(t, 404, w.Code)
}
