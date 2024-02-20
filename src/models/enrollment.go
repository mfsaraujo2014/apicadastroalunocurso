package models

import "errors"

type Enrollment struct {
	Code        uint64 `json:"code,omitempty"`
	StudentCode uint64 `json:"studentcode,omitempty"`
	CourseCode  uint64 `json:"coursecode,omitempty"`
}

func (enrollment *Enrollment) Prepare() error {
	if erro := enrollment.validate(); erro != nil {
		return erro
	}
	return nil
}

func (enrollment *Enrollment) validate() error {
	if enrollment.StudentCode == 0 {
		return errors.New("O Codigo de estudante e obrigatorio e nao pode estar em branco")
	}
	if enrollment.CourseCode == 0 {
		return errors.New("O Codigo do curso e obrigatorio e nao pode estar em branco")
	}
	return nil
}
