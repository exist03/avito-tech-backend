package internal

import "time"

type User struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type Segment struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type UserSegment struct {
	Segment `json:"segment"`
	TTL     time.Time `json:"ttl,omitempty"`
}
type UpdateRequest struct {
	UserId      int           `json:"user_id,omitempty"`
	SegmentsAdd []UserSegment `json:"segments_add,omitempty"`
	SegmentsDel []UserSegment `json:"segments_del,omitempty"`
}
