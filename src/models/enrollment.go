package models

type Enrollment struct {
	Code        uint64 `json:"code,omitempty"`
	StudentCode uint64 `json:"studentcode,omitempty"`
	CourseCode  uint64 `json:"coursecode,omitempty"`
}
