package constants

type SearchOrder int

const (
	SearchOrderDescending SearchOrder = 0
	SearchOrderAscending  SearchOrder = 1
)

var SearchOrders = []SearchOrder{
	SearchOrderDescending,
	SearchOrderAscending,
}

func (order SearchOrder) Values() []SearchOrder {
	return SearchOrders
}

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
	BeatmapSortUpdated    BeatmapSort = 9
)

var BeatmapSortingOptions = []BeatmapSort{
	BeatmapSortTitle,
	BeatmapSortArtist,
	BeatmapSortCreator,
	BeatmapSortDifficulty,
	BeatmapSortRanked,
	BeatmapSortRating,
	BeatmapSortPlays,
	BeatmapSortCreated,
	BeatmapSortRelevance,
	BeatmapSortUpdated,
}

func (sort BeatmapSort) Values() []BeatmapSort {
	return BeatmapSortingOptions
}

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

var BeatmapCategories = []BeatmapCategory{
	BeatmapCategoryAny,
	BeatmapCategoryLeaderboard,
	BeatmapCategoryRanked,
	BeatmapCategoryQualified,
	BeatmapCategoryLoved,
	BeatmapCategoryApproved,
	BeatmapCategoryPending,
	BeatmapCategoryWIP,
	BeatmapCategoryGraveyard,
}

func (category BeatmapCategory) Values() []BeatmapCategory {
	return BeatmapCategories
}
