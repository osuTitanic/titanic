# Location

This module provides geolocation lookups, resolving an IP address into a `Location` struct. The default `Provider` resolves geolocation data through a rest api backed by [ip-api.com](https://ip-api.com).

## Usage

Create the default provider, call `Setup()`, then `Resolve` an address.

```go
provider := location.NewProvider()
if err := provider.Setup(); err != nil {
	return err
}
// You may also use `state.Location` if you have the app set up

loc, err := provider.Resolve("1.1.1.1")
if err != nil {
	return err
}
fmt.Println(loc.CountryName, loc.City)
```

Local / Non-routable addresses are resolved against the host's own public address. On failure `Resolve` returns a `DefaultLocation` alongside the error, so callers can continue to use the location struct.
