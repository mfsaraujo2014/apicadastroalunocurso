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

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) CreateCourse(course models.Course) (int64, error) {
	statement, err := r.db.Prepare(
		"INSERT INTO curso(descricao, ementa) VALUES($1, $2) RETURNING codigo",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	var lastInsertID int64
	err = statement.QueryRow(utils.RemoveAccents(course.Description), utils.RemoveAccents(course.CourseProgram)).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func filterCourses(filters []models.Filter, query string, params []interface{}) (string, []interface{}) {
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
		containsDescriptionFilter := contains(filters, "descricao")
		containsCourseProgramFilter := contains(filters, "ementa")

		whereConditions := []string{}

		if containsDescriptionFilter {
			for _, filter := range filters {
				if filter.Key == "descricao" {
					value := strings.ToLower(string(filter.Value))
					parts := strings.Fields(value)
					for _, part := range parts {
						params = append(params, "%"+utils.RemoveAccents(part)+"%")
						whereConditions = append(whereConditions, " LOWER(descricao) LIKE $"+strconv.Itoa(len(params)))
					}
					break
				}
			}
		}
		if containsCourseProgramFilter {
			for _, filter := range filters {
				if filter.Key == "ementa" {
					value := strings.ToLower(string(filter.Value))
					parts := strings.Fields(value)
					for _, part := range parts {
						params = append(params, "%"+utils.RemoveAccents(part)+"%")
						whereConditions = append(whereConditions, " LOWER(ementa) LIKE $"+strconv.Itoa(len(params)))
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

	query += " ORDER BY codigo"

	return query, params
}

func (r *CourseRepository) GetCourses(ctx context.Context, skip, take int64, filters []models.Filter) ([]models.Course, error) {
	params := []interface{}{}

	query := `SELECT codigo, descricao, ementa FROM curso`

	query, params = filterCourses(filters, query, params)

	if skip != 0 || take != 0 {
		query += ` LIMIT $1 OFFSET $2;`
		params = append(params, take, skip)
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

	courses := []models.Course{}

	for rows.Next() {
		course := models.Course{}

		if err := rows.Scan(
			&course.Code,
			&course.Description,
			&course.CourseProgram,
		); err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *CourseRepository) UpdateCourse(ctx context.Context, ID uint64, course models.Course) (int64, error) {
	statement, err := r.db.PrepareContext(ctx,
		"UPDATE curso SET descricao = $1, ementa = $2 WHERE codigo = $3",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.ExecContext(ctx, course.Description, course.CourseProgram, ID)
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

func (r *CourseRepository) GetCourseByID(ctx context.Context, id uint64) (models.Course, error) {
	var course models.Course

	query := `
		SELECT c.codigo, c.descricao, c.ementa, a.codigo, a.nome
		FROM curso c
		LEFT JOIN curso_aluno ac ON c.codigo = ac.codigo_curso
		LEFT JOIN aluno a ON ac.codigo_aluno = a.codigo
		WHERE c.codigo = $1
	`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Course{}, fmt.Errorf("curso com o ID %d não encontrado", id)
		}
		return models.Course{}, err
	}
	defer rows.Close()

	var students []models.Student

	for rows.Next() {
		var student models.Student
		var studentCode sql.NullInt64
		var studentName sql.NullString

		err := rows.Scan(&course.Code, &course.Description, &course.CourseProgram, &studentCode, &studentName)
		if err != nil {
			return models.Course{}, err
		}

		if studentCode.Valid && studentName.Valid {
			student.Code = studentCode.Int64
			student.Name = studentName.String
			students = append(students, student)
		}
	}
	if err := rows.Err(); err != nil {
		return models.Course{}, err
	}

	course.Students = students

	return course, nil
}

func (r *CourseRepository) GetCoursesByStudent(ctx context.Context, studentCode int64) ([]models.Course, error) {
	query := `
        SELECT
			c.codigo,
            c.descricao,
			c.ementa
        FROM
            curso c
        JOIN
			curso_aluno cs ON c.codigo = cs.codigo_curso
        WHERE
            cs.codigo_aluno = $1
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, studentCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course

	for rows.Next() {
		var course models.Course

		if err := rows.Scan(
			&course.Code,
			&course.Description,
			&course.CourseProgram,
		); err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}
