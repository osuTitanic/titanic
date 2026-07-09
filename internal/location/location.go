package location

import (
	"github.com/osuTitanic/titanic/internal/constants"
)

// Location describes the geolocation of an IP address.
type Location struct {
	IP           string
	Latitude     float64
	Longitude    float64
	CountryCode  string
	CountryName  string
	CountryIndex int
	Timezone     string
	City         string
}

// DefaultLocation returns a Location with default / unknown values.
func DefaultLocation() *Location {
	return &Location{
		IP:           "127.0.0.1",
		Latitude:     0.0,
		Longitude:    0.0,
		CountryCode:  "XX",
		CountryName:  "Unknown",
		CountryIndex: 0,
		Timezone:     "UTC",
		City:         "Unknown",
	}
}

// IsLocal reports whether the location's IP is a local address.
func (l *Location) IsLocal() bool {
	return IsLocalIP(l.IP)
}

// SetCountryCode assigns the country fields from an ISO country code.
func (l *Location) SetCountryCode(code string) {
	name := constants.GetCountryNameFromCode(code)
	if name == "" {
		l.CountryCode = "XX"
		l.CountryName = "Unknown"
		l.CountryIndex = 0
		return
	}

	l.CountryCode = code
	l.CountryName = name
	l.CountryIndex = int(constants.GetCountryIndexFromCode(code))
}
