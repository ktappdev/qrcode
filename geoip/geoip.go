package geoip

import (
	"fmt"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

type GeoIP struct {
	db *geoip2.Reader
}

func New(databasePath string) (*GeoIP, error) {
	dbBytes, err := os.ReadFile(databasePath)
	if err != nil {
		return nil, err
	}

	db, err := geoip2.FromBytes(dbBytes)
	if err != nil {
		return nil, err
	}

	return &GeoIP{db: db}, nil
}

func (g *GeoIP) Close() error {
	return g.db.Close()
}

func (g *GeoIP) LookupCity(ipStr string) (*geoip2.City, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	city, err := g.db.City(ip)
	if err != nil {
		return nil, err
	}

	return city, nil
}
