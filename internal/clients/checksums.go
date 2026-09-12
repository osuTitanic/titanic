package clients

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/state"
)

func IsValidChecksum(app *state.State, version Version, checksum string) (bool, error) {
	// We have a database of all osu! releases from 2014
	// (when the new updater was introduced) up until now
	official, err := IsOfficialRelease(app, checksum)
	if err != nil || official {
		return official, err
	}

	// If we are still dealing with an official osu! version
	// we want to check our own patched osu! clients that we
	// distribute on titanic's downloads page
	if isOfficialIdentifier(version.Identifier()) {
		return IsTitanicRelease(app, version, checksum)
	}

	// Finally, there's official support for modded osu! clients
	// Each of them gets their own identifier, e.g. "digital", "flandremod", ...
	return IsModdedRelease(app, version, checksum)
}

func IsOfficialRelease(app *state.State, checksum string) (bool, error) {
	file, err := app.ReleasesOfficial.FetchFileByChecksum(checksum)
	if err != nil {
		return false, fmt.Errorf("find official client checksum: %w", err)
	}
	return file != nil, nil
}

func IsTitanicRelease(app *state.State, version Version, checksum string) (bool, error) {
	valid, err := app.ReleasesTitanic.Exists(version.Date, checksum)
	if err != nil {
		return false, fmt.Errorf("find titanic client checksum: %w", err)
	}
	return valid, nil
}

func IsModdedRelease(app *state.State, version Version, checksum string) (bool, error) {
	valid, err := app.ReleasesModded.EntryExists(version.Identifier(), checksum)
	if err != nil {
		return false, fmt.Errorf("find modded client checksum: %w", err)
	}
	return valid, nil
}

func isOfficialIdentifier(identifier string) bool {
	switch identifier {
	case "stable", "stable40", "test", "tourney", "cuttingedge", "beta",
		"ubertest", "public_test", "ce45", "peppy", "dev",
		"arcade", "noxna", "fallback", "a", "b", "c", "d", "e":
		return true
	default:
		return false
	}
}
