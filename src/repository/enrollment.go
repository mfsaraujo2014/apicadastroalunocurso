package repository

import (
	"context"
	"database/sql"

	"github.com/mfsaraujo2014/apicadastroalunocurso/src/models"
)

type EnrollmentRepository struct {
	db *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) EnrollStudent(ctx context.Context, enrollment models.Enrollment) (int64, error) {
	statement, err := r.db.PrepareContext(ctx,
		"INSERT INTO curso_aluno (codigo_aluno, codigo_curso) VALUES ($1, $2) RETURNING codigo",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	var lastInsertID int64
	err = statement.QueryRowContext(ctx, enrollment.StudentCode, enrollment.CourseCode).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}
