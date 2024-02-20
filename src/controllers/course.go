package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mfsaraujo2014/apicadastroalunocurso/src/answers"
	"github.com/mfsaraujo2014/apicadastroalunocurso/src/models"
	"github.com/mfsaraujo2014/apicadastroalunocurso/src/repository"
)

func CreateCourse(courseRepo *repository.CourseRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var course models.Course

		err := json.NewDecoder(r.Body).Decode(&course)
		if err != nil {
			answers.Erro(w, http.StatusBadRequest, err)
			return
		}

		if err := course.Prepare(); err != nil {
			answers.Erro(w, http.StatusBadRequest, err)
			return
		}

		id, err := courseRepo.CreateCourse(course)
		if err != nil {
			answers.Erro(w, http.StatusInternalServerError, err)
			return
		}

		fmt.Printf("Curso criado com sucesso! ID: %d\n", id)

		answers.JSON(w, http.StatusCreated, course)
	}
}

func GetCourses(courseRepo *repository.CourseRepository, studentRepo *repository.StudentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := r.URL.Query()
		searchTerm := query.Get("descricao")

		var filters []models.Filter
		if searchTerm != "" {
			filter := models.Filter{
				Key:   "descricao",
				Value: searchTerm,
			}
			filters = append(filters, filter)
		}

		courses, err := courseRepo.GetCourses(ctx, 0, 0, filters)
		if err != nil {
			answers.Erro(w, http.StatusInternalServerError, err)
			return
		}

		var wg sync.WaitGroup
		wg.Add(len(courses))
		for i := range courses {
			go func(i int) {
				defer wg.Done()
				students, err := studentRepo.GetStudentsByCourse(ctx, courses[i].Code)
				if err != nil {
					answers.Erro(w, http.StatusInternalServerError, err)
					return
				}
				courses[i].Students = students
			}(i)
		}
		wg.Wait()

		answers.JSON(w, http.StatusOK, courses)
	}
}

func UpdateCourse(courseRepo *repository.CourseRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		parametros := mux.Vars(r)
		courseID, erro := strconv.ParseUint(parametros["courseID"], 10, 64)
		if erro != nil {
			answers.Erro(w, http.StatusBadRequest, erro)
			return
		}

		var course models.Course

		err := json.NewDecoder(r.Body).Decode(&course)
		if err != nil {
			answers.Erro(w, http.StatusBadRequest, err)
			return
		}

		if erro = course.Prepare(); erro != nil {
			answers.Erro(w, http.StatusBadRequest, erro)
			return
		}

		id, err := courseRepo.UpdateCourse(ctx, courseID, course)
		if err != nil {
			answers.Erro(w, http.StatusInternalServerError, err)
			return
		}

		fmt.Printf("Curso editado com sucesso! ID: %d\n", id)

		answers.JSON(w, http.StatusNoContent, nil)
	}
}

func GetCourseByID(courseRepo *repository.CourseRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		parametros := mux.Vars(r)

		courseID, erro := strconv.ParseUint(parametros["courseID"], 10, 64)
		if erro != nil {
			answers.Erro(w, http.StatusBadRequest, erro)
			return
		}

		course, erro := courseRepo.GetCourseByID(ctx, courseID)
		if erro != nil {
			answers.Erro(w, http.StatusInternalServerError, erro)
			return
		}

		answers.JSON(w, http.StatusOK, course)
	}
}
