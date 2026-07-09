package constants

type SearchOrder int

const (
	SearchOrderDescending SearchOrder = 0
	SearchOrderAscending  SearchOrder = 1
)

type BeatmapSort int

const (
	BeatmapSortTitle      BeatmapSort = 0
	BeatmapSortArtist     BeatmapSort = 1
	BeatmapSortCreator    BeatmapSort = 2
	BeatmapSortDifficulty BeatmapSort = 3
	BeatmapSortRanked     BeatmapSort = 4
	BeatmapSortRating     BeatmapSort = 5
	BeatmapSortPlays      BeatmapSort = 6
	BeatmapSortCreated    BeatmapSort = 7
	BeatmapSortRelevance  BeatmapSort = 8
)

type BeatmapCategory int

const (
	BeatmapCategoryAny         BeatmapCategory = 0
	BeatmapCategoryLeaderboard BeatmapCategory = 1
	BeatmapCategoryRanked      BeatmapCategory = 2
	BeatmapCategoryQualified   BeatmapCategory = 3
	BeatmapCategoryLoved       BeatmapCategory = 4
	BeatmapCategoryApproved    BeatmapCategory = 5
	BeatmapCategoryPending     BeatmapCategory = 6
	BeatmapCategoryWIP         BeatmapCategory = 7
	BeatmapCategoryGraveyard   BeatmapCategory = 8
)
