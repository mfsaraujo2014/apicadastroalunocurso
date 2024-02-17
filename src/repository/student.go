package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mfsaraujo2014/apicadastroalunocurso/src/models"
	"github.com/mfsaraujo2014/apicadastroalunocurso/utils"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) CreateStudent(ctx context.Context, student models.Student) (int64, error) {
	statement, err := r.db.PrepareContext(ctx,
		"INSERT INTO aluno(nome) VALUES($1) RETURNING codigo",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	var lastInsertID int64
	err = statement.QueryRowContext(ctx, utils.RemoveAccents(student.Name)).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func filterStudents(filters []models.Filter, query string, params []interface{}) (string, []interface{}) {
	params = []interface{}{}
	contains := func(haystack []models.Filter, needle string) bool {
		for _, val := range haystack {
			if val.Key == needle {
				return true
			}
		}
		return false
	}

	if len(filters) > 0 {
		containsNameFilter := contains(filters, "nome")

		whereConditions := []string{}

		if containsNameFilter {
			for _, filter := range filters {
				if filter.Key == "nome" {
					value := strings.ToLower(string(filter.Value))
					parts := strings.Fields(value)
					for _, part := range parts {
						params = append(params, "%"+utils.RemoveAccents(part)+"%")
						whereConditions = append(whereConditions, " LOWER(nome) LIKE $"+strconv.Itoa(len(params)))
					}
					break
				}
			}
		}

		whereClause := strings.Join(whereConditions, " AND ")
		if len(whereClause) > 0 {
			if strings.Contains(query, "WHERE") {
				query += " AND (" + whereClause + ")"
			} else {
				query += " WHERE (" + whereClause + ")"
			}
		}
	}

	return query, params
}

func (r *StudentRepository) GetStudents(ctx context.Context, skip, take int64, filters []models.Filter) ([]models.Student, error) {
	params := []interface{}{}

	query := `SELECT codigo, nome FROM aluno`

	query, params = filterStudents(filters, query, params)

	fmt.Println(query)

	if skip != 0 || take != 0 {
		query += ` LIMIT $1 OFFSET $2;`
		params = append(params, take, skip)
	}

	for _, param := range params {
		fmt.Println("params: ", param)
	}

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := []models.Student{}

	for rows.Next() {
		student := models.Student{}

		if err := rows.Scan(
			&student.Code,
			&student.Name,
		); err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (r *StudentRepository) UpdateStudent(ctx context.Context, ID uint64, student models.Student) (int64, error) {
	statement, err := r.db.PrepareContext(ctx,
		"UPDATE aluno SET nome = $1 WHERE codigo = $2",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.ExecContext(ctx, student.Name, ID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, errors.New("nenhum aluno atualizado")
	}

	return int64(ID), nil
}

func (r *StudentRepository) GetStudentByID(ctx context.Context, id uint64) (models.Student, error) {
	var student models.Student

	query := `
		SELECT a.codigo, a.nome, c.codigo, c.descricao, c.ementa
		FROM aluno a
		LEFT JOIN curso_aluno ac ON a.codigo = ac.codigo_aluno
		LEFT JOIN curso c ON ac.codigo_curso = c.codigo
		WHERE a.codigo = $1
	`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Student{}, fmt.Errorf("estudante com o ID %d não encontrado", id)
		}
		return models.Student{}, err
	}
	defer rows.Close()

	var courses []models.Course

	for rows.Next() {
		var course models.Course
		var courseCode sql.NullInt64
		var courseDescription sql.NullString
		var courseCourseProgram sql.NullString
		err := rows.Scan(&student.Code, &student.Name, &courseCode, &courseDescription, &courseCourseProgram)
		if err != nil {
			return models.Student{}, err
		}
		if courseCode.Valid && courseDescription.Valid && courseCourseProgram.Valid {
			course.Code = courseCode.Int64
			course.Description = courseDescription.String
			course.CourseProgram = courseCourseProgram.String
			courses = append(courses, course)
		}
	}
	if err := rows.Err(); err != nil {
		return models.Student{}, err
	}

	student.Courses = courses

	return student, nil
}
