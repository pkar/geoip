package geoip

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

type locationHandler struct {
	LocationMap *LocationMap
	BlocksMap   *BlocksMap
}

// Run will initialize a database from flag options and begin a
// server for the api.
func Run(listen, locationPath, blocksPath string) error {
	f, err := os.Open(locationPath)
	if err != nil {
		return err
	}
	lm, err := NewLocationMap(f)
	if err != nil {
		f.Close()
		return err
	}
	f.Close()

	f, err = os.Open(blocksPath)
	if err != nil {
		return err
	}
	bm, err := NewBlocksMap(f)
	if err != nil {
		f.Close()
		return err
	}
	f.Close()
	locHandler := &locationHandler{
		LocationMap: lm,
		BlocksMap:   bm,
	}

	mux := http.NewServeMux()
	mux.Handle("/location", locHandler)
	mux.Handle("/location/", locHandler)
	log.Println("listening on", listen)
	return http.ListenAndServe(listen, mux)
}

func (lh *locationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	ip := vals.Get("ip")
	locID := vals.Get("locId")
	switch {
	case ip != "":
		id, ok := lh.BlocksMap.Lookup(ip)
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		l, ok := lh.LocationMap.LookupByID(id)
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(l.Serialized)
		return
	case locID != "":
		id, err := strconv.Atoi(locID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+" invalid locId", http.StatusBadRequest)
			return
		}
		l, ok := lh.LocationMap.LookupByID(id)
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(l.Serialized)
		return
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
