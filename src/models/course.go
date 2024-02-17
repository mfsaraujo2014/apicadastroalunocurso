package models

import (
	"errors"
)

type Course struct {
	Code          int64     `json:"code,omitempty"`
	Description   string    `json:"description,omitempty"`
	CourseProgram string    `json:"courseprogram,omitempty"`
	Students      []Student `json:"students,omitempty"`
}

func (course *Course) Prepare() error {
	if erro := course.validate(); erro != nil {
		return erro
	}
	return nil
}

func (course *Course) validate() error {
	if course.Description == "" {
		return errors.New("A Descricao é obrigatória e não pode estar em branco")
	}
	return nil
}
