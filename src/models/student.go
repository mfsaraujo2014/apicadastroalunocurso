package models

import "errors"

type Student struct {
	Code    int64    `json:"code,omitempty"`
	Name    string   `json:"name,omitempty"`
	Courses []Course `json:"courses,omitempty"`
}

func (student *Student) Prepare() error {
	if erro := student.validate(); erro != nil {
		return erro
	}
	return nil
}

func (student *Student) validate() error {
	if student.Name == "" {
		return errors.New("O Nome é obrigatório e não pode estar em branco")
	}
	return nil
}
