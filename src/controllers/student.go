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

func CreateStudent(studentRepo *repository.StudentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var student models.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			answers.Erro(w, http.StatusBadRequest, err)
			return
		}

		if err := student.Prepare(); err != nil {
			answers.Erro(w, http.StatusBadRequest, err)
			return
		}

		id, err := studentRepo.CreateStudent(ctx, student)
		if err != nil {
			answers.Erro(w, http.StatusInternalServerError, err)
			return
		}

		fmt.Printf("Estudante criado com sucesso! ID: %d\n", id)

		answers.JSON(w, http.StatusCreated, student)
	}
}

func GetStudents(studentRepo *repository.StudentRepository, courseRepo *repository.CourseRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := r.URL.Query()
		fmt.Println(query)
		searchTerm := query.Get("nome")
		fmt.Println("searchTerm: ", searchTerm)

		var filters []models.Filter
		if searchTerm != "" {
			filter := models.Filter{
				Key:   "nome",
				Value: searchTerm,
			}
			filters = append(filters, filter)
		}

		students, err := studentRepo.GetStudents(ctx, 0, 0, filters)
		if err != nil {
			answers.Erro(w, http.StatusInternalServerError, err)
			return
		}

		var wg sync.WaitGroup
		wg.Add(len(students))
		for i := range students {
			go func(i int) {
				defer wg.Done()
				courses, err := courseRepo.GetCoursesByStudent(ctx, students[i].Code)
				if err != nil {
					answers.Erro(w, http.StatusInternalServerError, err)
					return
				}
				students[i].Courses = courses
			}(i)
		}
		wg.Wait()

		answers.JSON(w, http.StatusOK, students)
	}
}

func UpdateStudent(studentRepo *repository.StudentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		parametros := mux.Vars(r)
		studentID, erro := strconv.ParseUint(parametros["studentID"], 10, 64)
		if erro != nil {
			answers.Erro(w, http.StatusBadRequest, erro)
			return
		}

		var student models.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			answers.Erro(w, http.StatusBadRequest, err)
			return
		}

		if erro = student.Prepare(); erro != nil {
			answers.Erro(w, http.StatusBadRequest, erro)
			return
		}

		id, err := studentRepo.UpdateStudent(ctx, studentID, student)
		if err != nil {
			answers.Erro(w, http.StatusInternalServerError, err)
			return
		}

		fmt.Printf("Estudante editado com sucesso! ID: %d\n", id)

		answers.JSON(w, http.StatusNoContent, nil)
	}
}

// for _, b := range bookings {
// 	g, ctx := errgroup.WithContext(ctx)

// 	g.Go(func() error {
// 		b.SalesOrder, err = s.salesorder.Find(ctx, b.SalesOrder.ID, requestCompany)
// 		return err
// 	})

// 	if err := g.Wait(); err != nil {
// 		return nil, err
// 	}
// }

func GetStudentByID(studentRepo *repository.StudentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		parametros := mux.Vars(r)

		studentID, erro := strconv.ParseUint(parametros["studentID"], 10, 64)
		if erro != nil {
			answers.Erro(w, http.StatusBadRequest, erro)
			return
		}

		student, erro := studentRepo.GetStudentByID(ctx, studentID)
		if erro != nil {
			answers.Erro(w, http.StatusInternalServerError, erro)
			return
		}

		answers.JSON(w, http.StatusOK, student)
	}
}
