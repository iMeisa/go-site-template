package repository

import (
	"github.com/iMeisa/errortrace"
	"github.com/iMeisa/go-site-template/internal/models"
)

type DatabaseRepo interface {

	// Query functions
	QueryUsers(params map[string]string) ([]models.User, errortrace.ErrorTrace)

	// User functions
	InsertUser(user models.User) errortrace.ErrorTrace
	LoginUpdateUser(user *models.User) errortrace.ErrorTrace
	UpdateUser(user models.User) errortrace.ErrorTrace
}
