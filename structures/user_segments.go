package structures

type SegmentUsers struct {
	Slug    string   `json:"slug" binding:"required"`
	UserIDs []string `json:"user_ids" binding:"required"`
}
