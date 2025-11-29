package geoip

import (
	"fmt"
	"net/netip"

	"github.com/oschwald/maxminddb-golang/v2"
)

type Locator interface {
	Lookup(ip netip.Addr) (*Location, error)
	Close() error
}

type Location struct {
	Country    string
	CountryISO string
	City       string
}

type maxmindLocator struct {
	db *maxminddb.Reader
}

func NewLocator() (Locator, error) {
	db, err := maxminddb.Open("./GeoLite2-City.mmdb")
	if err != nil {
		return nil, fmt.Errorf("open maxmind: %w", err)
	}
	return &maxmindLocator{db: db}, nil
}

func (m *maxmindLocator) Lookup(ip netip.Addr) (*Location, error) {
	var rec struct {
		Country struct {
			ISOCode string            `maxminddb:"iso_code"`
			Names   map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
		City struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"city"`
	}

	if err := m.db.Lookup(ip).Decode(&rec); err != nil {
		return nil, err
	}
	return &Location{
		Country:    rec.Country.Names["en"],
		CountryISO: rec.Country.ISOCode,
		City:       rec.City.Names["en"],
	}, nil
}

func (m *maxmindLocator) Close() error { return m.db.Close() }
