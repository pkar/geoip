package geoip

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

var (
	// ErrInvalidIP is for invalid ip addresses.
	ErrInvalidIP = errors.New("invalid ip provided")
)

func ipToInt(ip string) (int, error) {
	ints := strings.Split(ip, ".")
	if len(ints) != 4 {
		return 0, ErrInvalidIP
	}
	i1, err := strconv.Atoi(ints[0])
	if err != nil {
		log.Println(err)
		return 0, ErrInvalidIP
	}
	i2, err := strconv.Atoi(ints[1])
	if err != nil {
		log.Println(err)
		return 0, ErrInvalidIP
	}
	i3, err := strconv.Atoi(ints[2])
	if err != nil {
		log.Println(err)
		return 0, ErrInvalidIP
	}
	i4, err := strconv.Atoi(ints[3])
	if err != nil {
		log.Println(err)
		return 0, ErrInvalidIP
	}
	return (16777216 * i1) + (65536 * i2) + (256 * i3) + i4, nil
}

func intToIP(i int) (string, error) {
	i1 := strconv.Itoa(int(i/16777216) % 256)
	i2 := strconv.Itoa(int(i/65536) % 256)
	i3 := strconv.Itoa(int(i/256) % 256)
	i4 := strconv.Itoa(int(i) % 256)
	ips := []string{i1, i2, i3, i4}
	return strings.Join(ips, "."), nil
}
