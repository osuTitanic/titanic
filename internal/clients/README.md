# Clients

This module contains helpers for validating osu! clients, i.e. version parsing & checksum validations.
It supports official osu! releases, Titanic's patched clients & modded clients.

## Version parsing

Use `ParseVersion` to parse the version string reported by an osu! client.

```go
version, ok := clients.ParseVersion(rawVersion)
if !ok {
	return errors.New("invalid client version")
}

fmt.Println(version.Date, version.Revision, version.Identifier())
```

## Checksum validation

Use `IsValidChecksum` with the parsed version and checksum reported by the client.

```go
valid, err := clients.IsValidChecksum(app, version, checksum)
if err != nil {
	return err
}
if !valid {
	return errors.New("invalid client checksum")
}
```

It checks the known official osu! release files first. If no file matches, client identifiers such as "stable", "cuttingedge", "beta", etc. are checked against Titanic's patched releases. Other identifiers are checked against the corresponding modded client.

The individual `IsOfficialRelease`, `IsTitanicRelease` & `IsModdedRelease` helpers are also available when one would need to validate them directly.
