package rankings

import "errors"

var ErrStatsIsNil = errors.New("rankings: stats is nil")
var ErrNoPlayerAbove = errors.New("rankings: no player above")
var ErrRedisClientNotInitialized = errors.New("rankings: redis client is not initialized")
