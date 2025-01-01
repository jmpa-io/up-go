package up

type CategoryAttributes struct {
	Name string `json:"name"`
}

type CategoryRelationships struct {
	Parent struct {
	} `json:"parent"`
	Child struct {
	} `json:"child"`
	Links Links `json:"links"`
}

// CategoryData represents a category in Up.
type CategoryData Data[CategoryAttributes, CategoryRelationships]

// CategoryPaginationWrapper is a pagination wrapper for a slice of CategoryData
type CategoryPaginationWrapper WrapperSlice[CategoryData]
