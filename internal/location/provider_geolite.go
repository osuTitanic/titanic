package location

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
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
	_, err := os.Stat(provider.databasePath)
	wasNotFound := errors.Is(err, os.ErrNotExist)

	if err != nil && !wasNotFound {
		return err
	}
	if wasNotFound {
		if err := provider.Download(); err != nil {
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

	if location.CountryCode == "" || location.CountryCode == "XX" {
		return location, errors.New("missing country information")
	}
	return location, nil
}

func (provider *GeoLiteProvider) Download() error {
	lock, err := provider.lockOrWait(provider.databasePath + ".lock")
	if err != nil {
		return fmt.Errorf("lock database: %w", err)
	}
	defer lock.Close()

	// Another service may have finished the download while this one waited for the lock
	if _, err := os.Stat(provider.databasePath); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	provider.logger.Info(
		"Downloading GeoLite2 database...",
		"url", provider.downloadUrl, "path", provider.databasePath,
	)

	// Create a temporary file in the same directory as the target database file
	// Once the download finished we can move it to the target path
	file, err := os.CreateTemp(
		filepath.Dir(provider.databasePath),
		"."+filepath.Base(provider.databasePath)+".tmp-*",
	)
	if err != nil {
		return err
	}

	temporaryPath := file.Name()
	defer func() {
		file.Close()
		os.Remove(temporaryPath)
	}()

	response, err := http.Get(provider.downloadUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("download database: unexpected status code: %d", response.StatusCode)
	}

	if _, err := io.Copy(file, response.Body); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Chmod(0o644); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}

	if err := os.Rename(temporaryPath, provider.databasePath); err != nil {
		return err
	}
	return nil
}

func (provider *GeoLiteProvider) lockOrWait(path string) (*flock.Flock, error) {
	fileLock := flock.New(path)

	// Try to acquire an exclusive lock on the file without blocking
	locked, err := fileLock.TryLock()
	if err != nil {
		return nil, err
	}
	if locked {
		return fileLock, nil
	}

	provider.logger.Info(
		"Waiting for GeoLite2 database download...",
		"path", provider.databasePath,
	)
	if err := fileLock.Lock(); err != nil {
		return nil, err
	}
	return fileLock, nil
}
