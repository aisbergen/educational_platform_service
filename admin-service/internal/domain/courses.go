package domain

import "time"

type Courses struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Students    []string  `json:"students"`
}

type CoursesOut struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Students    []string `json:"students"`
}

type ResponseCourse struct {
	Courses `json:"course"`
}
