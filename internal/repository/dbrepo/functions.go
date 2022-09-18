package dbrepo

import (
	"context"
	"github.com/iMeisa/errortrace"
	"github.com/iMeisa/go-site-template/internal/models"
	"time"
)

// InsertUser inserts a user into the database
func (m *dbRepo) InsertUser(user models.User) errortrace.ErrorTrace {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `
		INSERT INTO users (
				     user_id
				   , username
				   , discriminator
				   , pfp_hash
				   , free_months
		)  VALUES ($1, $2, $3, $4, $5)
	`
	_, err := m.DB.ExecContext(ctx, statement, user.ID, user.Username, user.Discriminator, user.PfpHash, 0)
	if err != nil {
		return errortrace.NewTrace(err)
	}

	return errortrace.NilTrace()
}

// LoginUpdateUser updates a user's data if changed
func (m *dbRepo) LoginUpdateUser(user *models.User) errortrace.ErrorTrace {

	// Query user
	query := EmptyQuery()
	query["user_id"] = user.ID
	users, trace := m.QueryUsers(query)
	if trace.HasError() {
		return trace
	}

	// If user doesn't exist, insert
	if len(users) < 1 {
		trace = m.InsertUser(*user)
		if trace.HasError() {
			return trace
		}
	}

	// Check if user is the same
	if users[0].Username == user.Username && users[0].Discriminator == user.Discriminator && users[0].PfpHash == user.PfpHash {
		return errortrace.NilTrace()
	}

	// Update user
	trace = m.UpdateUser(*user)
	if trace.HasError() {
		return trace
	}

	return errortrace.NilTrace()
}

// UpdateUser updates a user in the database
func (m *dbRepo) UpdateUser(user models.User) errortrace.ErrorTrace {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `
		UPDATE users 
		SET username=$1, discriminator=$2, pfp_hash=$3
		WHERE user_id=$4
 	`
	_, err := m.DB.ExecContext(ctx, statement, user.Username, user.Discriminator, user.PfpHash, user.ID)
	if err != nil {
		return errortrace.NewTrace(err)
	}

	return errortrace.NilTrace()
}
