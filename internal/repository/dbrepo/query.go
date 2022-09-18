package dbrepo

import (
	"context"
	"database/sql"
	"github.com/iMeisa/errortrace"
	"github.com/iMeisa/go-site-template/internal/models"
	"time"
)

func (m *dbRepo) QueryTableCols(schema, tableName string) (map[string]string, errortrace.ErrorTrace) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// VTP queries can be null
	statement := `
		SELECT column_name
		     , data_type
		FROM information_schema.columns
		WHERE table_schema=$1
		AND table_name=$2
	 `

	rows, err := m.DB.QueryContext(ctx, statement, schema, tableName)
	if err != nil {
		return nil, errortrace.NewTrace(err)
	}

	cols := make(map[string]string)

	for rows.Next() {
		var colName string
		var dataType string

		if err = rows.Scan(&colName, &dataType); err != nil {
			return nil, errortrace.NewTrace(err)
		}

		cols[colName] = dataType
	}

	return cols, errortrace.NilTrace()
}

// QueryUsers with params returns a list of users
func (m *dbRepo) QueryUsers(params map[string]string) ([]models.User, errortrace.ErrorTrace) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// php_hash can be null
	statement := `
		SELECT user_id
		     , username
		     , discriminator
		     , pfp_hash
		     , free_months
		FROM users
	`

	if len(params) > 0 {
		statement += `WHERE ` + createParamStatement(params, " AND ")
	}

	statement += ` ORDER BY username`

	rows, err := m.DB.QueryContext(ctx, statement)
	if err != nil {
		return nil, errortrace.NewTrace(err)
	}

	var users []models.User

	for rows.Next() {
		// User model
		var nextUser models.User
		// Nullable values
		var pfpHash sql.NullString

		// Same order as query
		err = rows.Scan(
			&nextUser.ID,
			&nextUser.Username,
			&nextUser.Discriminator,
			&pfpHash,
		)
		if err != nil {
			return nil, errortrace.NewTrace(err)
		}

		// Check if null
		if pfpHash.Valid {
			nextUser.PfpHash = pfpHash.String
		}

		users = append(users, nextUser)
	}

	return users, errortrace.NilTrace()
}
