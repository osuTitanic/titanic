package resources

import "errors"

var ErrNoMirrorsAvailable = errors.New("resources: no mirrors available")
var ErrResourceNotFound = errors.New("resources: resource not found")
var ErrRateLimited = errors.New("resources: mirror is rate limited")
