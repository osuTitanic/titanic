package location

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/netip"
	"os"

	"github.com/oschwald/geoip2-golang/v2"
)

type GeoLiteProvider struct {
	databasePath string
	downloadUrl  string
	logger       *slog.Logger
	database     *geoip2.Reader
}

func NewGeoLiteProvider(databasePath, downloadUrl string) *GeoLiteProvider {
	return &GeoLiteProvider{
		databasePath: databasePath,
		downloadUrl:  downloadUrl,
		logger:       slog.Default().With("component", "location", "provider", "geolite"),
	}
}

func (provider *GeoLiteProvider) Setup() error {
	if _, err := os.Stat(provider.databasePath); err != nil {
		err = provider.Download()
		if err != nil {
			return err
		}
	}

	db, err := geoip2.Open(provider.databasePath)
	if err == nil {
		// Database exists & is valid
		provider.database = db
		return nil
	}

	provider.logger.Error(
		"Invalid GeoLite2 database file",
		"error", err,
	)

	// Remove invalid database file & let it re-download on next startup
	err = os.Remove(provider.databasePath)
	if err != nil {
		return err
	}

	return errors.New("invalid GeoLite2 database file")
}

func (provider *GeoLiteProvider) Close() error {
	if provider.database == nil {
		return nil
	}

	err := provider.database.Close()
	provider.database = nil
	return err
}

func (provider *GeoLiteProvider) Resolve(ipString string) (*Location, error) {
	ip, err := netip.ParseAddr(ipString)
	if err != nil {
		return nil, err
	}

	if provider.database == nil {
		return nil, errors.New("GeoLite2 database is not initialized")
	}

	record, err := provider.database.City(ip)
	if err != nil {
		provider.logger.Error(
			"Failed to resolve IP address",
			"ip", ipString, "error", err,
		)
		return nil, err
	}
	if !record.HasData() {
		return nil, errors.New("no geolocation data found for IP address")
	}

	location := &Location{IP: ipString}

	if record.Location.Latitude != nil {
		location.Latitude = *record.Location.Latitude
	}
	if record.Location.Longitude != nil {
		location.Longitude = *record.Location.Longitude
	}

	if record.Location.TimeZone != "" {
		location.Timezone = record.Location.TimeZone
	}
	if record.City.Names.English != "" {
		location.City = record.City.Names.English
	}
	if record.Country.ISOCode != "" {
		location.SetCountryCode(record.Country.ISOCode)
	}

	return location, nil
}

func (provider *GeoLiteProvider) Download() error {
	provider.logger.Info(
		"Downloading GeoLite2 database...",
		"url", provider.downloadUrl, "path", provider.databasePath,
	)

	// Create a new database file
	file, err := os.Create(provider.databasePath)
	if err != nil {
		return err
	}
	defer file.Close()

	response, err := http.Get(provider.downloadUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
