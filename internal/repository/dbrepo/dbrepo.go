package dbrepo

import (
	"database/sql"

	"github.com/taldrori/bookings/internal/config"
	"github.com/taldrori/bookings/internal/repository"
)

type postgressDBRepo struct {
	App *config.Appconfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.Appconfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.Appconfig) repository.DatabaseRepo {
	return &postgressDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.Appconfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
