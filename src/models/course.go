package models

import (
	"errors"
)

type Course struct {
	Code          int64     `json:"code,omitempty"`
	Description   string    `json:"description,omitempty"`
	CourseProgram string    `json:"courseprogram,omitempty"`
	Students      []Student `json:"students"`
}

func (course *Course) Prepare() error {
	if erro := course.validate(); erro != nil {
		return erro
	}
	return nil
}

func (course *Course) validate() error {
	if course.Description == "" {
		return errors.New("A Descricao e obrigatoria e nao pode estar em branco")
	}
	if course.CourseProgram == "" {
		return errors.New("A Ementa e obrigatoria e nao pode estar em branco")
	}
	return nil
}
