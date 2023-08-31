package domain

import "time"

type User struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type Segment struct {
	Id      int       `json:"id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Percent float64   `json:"percent"`
	TTL     time.Time `json:"ttl,omitempty"`
}
type UpdateRequest struct {
	UserId      int       `json:"user_id,omitempty"`
	SegmentsAdd []Segment `json:"segments_add,omitempty"`
	SegmentsDel []Segment `json:"segments_del,omitempty"`
}
