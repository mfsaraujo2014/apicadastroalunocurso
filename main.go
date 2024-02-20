package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/mfsaraujo2014/apicadastroalunocurso/src/controllers"
	"github.com/mfsaraujo2014/apicadastroalunocurso/src/repository"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=9144 dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalf("db: failed to connect./n%s", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to get postgres driver instance: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	courseRepo := repository.NewCourseRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)

	r := mux.NewRouter()
	r.HandleFunc("/courses", controllers.CreateCourse(courseRepo)).Methods(http.MethodPost)
	r.HandleFunc("/courses", controllers.GetCourses(courseRepo, studentRepo)).Methods(http.MethodGet)
	r.HandleFunc("/courses/{courseID}", controllers.GetCourseByID(courseRepo)).Methods(http.MethodGet)
	r.HandleFunc("/courses/{courseID}", controllers.UpdateCourse(courseRepo)).Methods(http.MethodPut)
	r.HandleFunc("/students", controllers.CreateStudent(studentRepo)).Methods(http.MethodPost)
	r.HandleFunc("/students", controllers.GetStudents(studentRepo, courseRepo)).Methods(http.MethodGet)
	r.HandleFunc("/students/{studentID}", controllers.GetStudentByID(studentRepo)).Methods(http.MethodGet)
	r.HandleFunc("/students/{studentID}", controllers.UpdateStudent(studentRepo)).Methods(http.MethodPut)
	r.HandleFunc("/enrollments", controllers.EnrollStudent(courseRepo, studentRepo, enrollmentRepo)).Methods(http.MethodPost)

	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	fmt.Printf("Escutando na porta %d", 9001)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 9001), handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(r)))
}
