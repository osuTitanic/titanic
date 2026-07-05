package constants

import "regexp"

var (
	OsuVersion        = regexp.MustCompile(`^b(?P<date>\d{1,8})(?:(?P<name>\w+\b))?(?:\.(?P<revision>\d{1,2}|))?(?P<stream>\w+)?$`)
	OsuUserAgent      = regexp.MustCompile(`^osu!*|Mozilla/\d+\.\d+ \(compatible; Clever Internet Suite \d+\.\d+\)$`)
	OsuChatLinkModern = regexp.MustCompile(`\[((?:https?:\/\/)[^\s\]]+)\s+(.+?)\]`)
	OsuChatLinkLegacy = regexp.MustCompile(`\[([^\]]+)\]\((https?:\/\/[^)]+)\)`)
	Email             = regexp.MustCompile(`^[^@\s]{1,200}@[^@\s\.]{1,30}(?:\.[^@\.\s]{2,24})+$`)
	Username          = regexp.MustCompile(`^[a-zA-Z0-9^\-{}_\[\] ]+$`)
	DiscordUsername   = regexp.MustCompile(`^@?[a-z0-9_-]{3,32}$`)
	DiscordEmote      = regexp.MustCompile(`<a?:([a-zA-Z0-9_]{2,32}):\d{17,20}>`)
	TwitterHandle     = regexp.MustCompile(`https?://(www.)?(twitter|x)\.com/(@\w+|\w+)`)
	URL               = regexp.MustCompile(`(?i)\b((?:https?://|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}/)(?:[^\s()<>]+|\(([^\s()<>]+|(\([^\s()<>]+\)))*\))+(?:\(([^\s()<>]+|(\([^\s()<>]+\)))*\)|[^\s` + "`" + `!()\[\]{};:'\".,<>?«»“”‘’]))`) // yeahhh this is something alright
)
