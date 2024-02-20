package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfsaraujo2014/apicadastroalunocurso/src/answers"
	"github.com/mfsaraujo2014/apicadastroalunocurso/src/models"
	"github.com/mfsaraujo2014/apicadastroalunocurso/src/repository"
)

func EnrollStudent(courseRepo *repository.CourseRepository, studentRepo *repository.StudentRepository, enrollmentRepo *repository.EnrollmentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var enrollment models.Enrollment
		err := json.NewDecoder(r.Body).Decode(&enrollment)
		if err != nil {
			answers.Erro(w, http.StatusBadRequest, err)
			return
		}

		if erro := enrollment.Prepare(); erro != nil {
			answers.Erro(w, http.StatusBadRequest, erro)
			return
		}

		course, err := courseRepo.GetCourseByID(ctx, enrollment.CourseCode)
		if err != nil {
			answers.Erro(w, http.StatusInternalServerError, err)
			return
		}
		if len(course.Students) >= 10 {
			answers.Erro(w, http.StatusBadRequest, errors.New("O curso esta cheio"))
			return
		}

		student, err := studentRepo.GetStudentByID(ctx, enrollment.StudentCode)
		if err != nil {
			answers.Erro(w, http.StatusInternalServerError, err)
			return
		}
		if len(student.Courses) >= 3 {
			answers.Erro(w, http.StatusBadRequest, errors.New("O aluno ja esta matriculado em 3 cursos"))
			return
		}
		for _, course := range student.Courses {
			if course.Code == int64(enrollment.CourseCode) {
				answers.Erro(w, http.StatusBadRequest, errors.New("O aluno ja esta matriculado neste curso"))
				return
			}
		}

		id, err := enrollmentRepo.EnrollStudent(ctx, enrollment)
		if err != nil {
			answers.Erro(w, http.StatusInternalServerError, err)
			return
		}
		fmt.Printf("Matrícula de numero %d realizada com sucesso!\n", id)

		answers.JSON(w, http.StatusCreated, enrollment)
	}
}
