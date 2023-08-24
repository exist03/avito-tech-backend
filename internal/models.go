package internal

type User struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type Segment struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type SegmentList struct {
	Segments []string
}
