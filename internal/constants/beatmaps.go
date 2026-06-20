package constants

type BeatmapStatus int

const (
	BeatmapStatusInactive  BeatmapStatus = -3
	BeatmapStatusGraveyard BeatmapStatus = -2
	BeatmapStatusWIP       BeatmapStatus = -1
	BeatmapStatusPending   BeatmapStatus = 0
	BeatmapStatusRanked    BeatmapStatus = 1
	BeatmapStatusApproved  BeatmapStatus = 2
	BeatmapStatusQualified BeatmapStatus = 3
	BeatmapStatusLoved     BeatmapStatus = 4
)

func (status BeatmapStatus) String() string {
	switch status {
	case BeatmapStatusInactive:
		return "Inactive"
	case BeatmapStatusGraveyard:
		return "Graveyard"
	case BeatmapStatusWIP:
		return "WIP"
	case BeatmapStatusPending:
		return "Pending"
	case BeatmapStatusRanked:
		return "Ranked"
	case BeatmapStatusApproved:
		return "Approved"
	case BeatmapStatusQualified:
		return "Qualified"
	case BeatmapStatusLoved:
		return "Loved"
	default:
		return "Unknown"
	}
}

type BeatmapServer int

const (
	BeatmapServerBancho  BeatmapServer = 0
	BeatmapServerTitanic BeatmapServer = 1
)

type BeatmapResourceType int

const (
	BeatmapResourceTypeOsz        BeatmapResourceType = 0
	BeatmapResourceTypeOszNoVideo BeatmapResourceType = 1
	BeatmapResourceTypeBeatmap    BeatmapResourceType = 2
	BeatmapResourceTypeThumbnail  BeatmapResourceType = 3
	BeatmapResourceTypeBackground BeatmapResourceType = 4
	BeatmapResourceTypeAudio      BeatmapResourceType = 5
)

type BeatmapGenre int

const (
	BeatmapGenreAny         BeatmapGenre = 0
	BeatmapGenreUnspecified BeatmapGenre = 1
	BeatmapGenreVideoGame   BeatmapGenre = 2
	BeatmapGenreAnime       BeatmapGenre = 3
	BeatmapGenreRock        BeatmapGenre = 4
	BeatmapGenrePop         BeatmapGenre = 5
	BeatmapGenreOther       BeatmapGenre = 6
	BeatmapGenreNovelty     BeatmapGenre = 7
	BeatmapGenreHipHop      BeatmapGenre = 9
	BeatmapGenreElectronic  BeatmapGenre = 10
	BeatmapGenreMetal       BeatmapGenre = 11
	BeatmapGenreClassical   BeatmapGenre = 12
	BeatmapGenreFolk        BeatmapGenre = 13
	BeatmapGenreJazz        BeatmapGenre = 14
)

func (genre BeatmapGenre) Value() int {
	return int(genre)
}

func (genre BeatmapGenre) String() string {
	switch genre {
	case BeatmapGenreAny:
		return "Any"
	case BeatmapGenreUnspecified:
		return "Unspecified"
	case BeatmapGenreVideoGame:
		return "Video Game"
	case BeatmapGenreAnime:
		return "Anime"
	case BeatmapGenreRock:
		return "Rock"
	case BeatmapGenrePop:
		return "Pop"
	case BeatmapGenreOther:
		return "Other"
	case BeatmapGenreNovelty:
		return "Novelty"
	case BeatmapGenreHipHop:
		return "Hip Hop"
	case BeatmapGenreElectronic:
		return "Electronic"
	case BeatmapGenreMetal:
		return "Metal"
	case BeatmapGenreClassical:
		return "Classical"
	case BeatmapGenreFolk:
		return "Folk"
	case BeatmapGenreJazz:
		return "Jazz"
	default:
		return "Unknown"
	}
}

type BeatmapLanguage int

const (
	BeatmapLanguageAny          BeatmapLanguage = 0
	BeatmapLanguageUnspecified  BeatmapLanguage = 1
	BeatmapLanguageEnglish      BeatmapLanguage = 2
	BeatmapLanguageJapanese     BeatmapLanguage = 3
	BeatmapLanguageChinese      BeatmapLanguage = 4
	BeatmapLanguageInstrumental BeatmapLanguage = 5
	BeatmapLanguageKorean       BeatmapLanguage = 6
	BeatmapLanguageFrench       BeatmapLanguage = 7
	BeatmapLanguageGerman       BeatmapLanguage = 8
	BeatmapLanguageSwedish      BeatmapLanguage = 9
	BeatmapLanguageSpanish      BeatmapLanguage = 10
	BeatmapLanguageItalian      BeatmapLanguage = 11
	BeatmapLanguageRussian      BeatmapLanguage = 12
	BeatmapLanguagePolish       BeatmapLanguage = 13
	BeatmapLanguageOther        BeatmapLanguage = 14
)

func (language BeatmapLanguage) Value() int {
	return int(language)
}

func (language BeatmapLanguage) String() string {
	switch language {
	case BeatmapLanguageAny:
		return "Any"
	case BeatmapLanguageUnspecified:
		return "Unspecified"
	case BeatmapLanguageEnglish:
		return "English"
	case BeatmapLanguageJapanese:
		return "Japanese"
	case BeatmapLanguageChinese:
		return "Chinese"
	case BeatmapLanguageInstrumental:
		return "Instrumental"
	case BeatmapLanguageKorean:
		return "Korean"
	case BeatmapLanguageFrench:
		return "French"
	case BeatmapLanguageGerman:
		return "German"
	case BeatmapLanguageSwedish:
		return "Swedish"
	case BeatmapLanguageSpanish:
		return "Spanish"
	case BeatmapLanguageItalian:
		return "Italian"
	case BeatmapLanguageRussian:
		return "Russian"
	case BeatmapLanguagePolish:
		return "Polish"
	case BeatmapLanguageOther:
		return "Other"
	default:
		return "Unknown"
	}
}
