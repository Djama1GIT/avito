package structures

type UserSegments struct {
	UserId           int      `json:"user_id" binding:"required"`
	SegmentsToAdd    []string `json:"segments_to_add" binding:"required"`
	SegmentsToDelete []string `json:"segments_to_delete" binding:"required"`
}
