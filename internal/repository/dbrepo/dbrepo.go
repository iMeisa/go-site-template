package dbrepo

import (
	"database/sql"
	"fmt"
	"github.com/iMeisa/go-site-template/internal/config"
	"github.com/iMeisa/go-site-template/internal/repository"
	"strconv"
	"strings"
)

type dbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewDBRepo(conn *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &dbRepo{
		App: app,
		DB:  conn,
	}
}

func EmptyQuery() map[string]string {
	return make(map[string]string)
}

func createParamStatement(params map[string]string, sep string) string {
	var paramsList []string

	for key, val := range params {
		param := fmt.Sprintf("%v='%v'", key, val)

		// If param is int, remove quotes from value
		if valNum, err := strconv.Atoi(val); err == nil && valNum < (1<<31-1) {
			param = fmt.Sprintf("%v=%v", key, val)
		}

		paramsList = append(paramsList, param)
	}

	return strings.Join(paramsList, sep)
}
