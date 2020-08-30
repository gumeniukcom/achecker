package dao

import (
	"context"
	"github.com/gumeniukcom/achecker/postgres"
	"github.com/gumeniukcom/achecker/resultdaemon/structs"
	"github.com/rs/zerolog/log"
	"time"
)

// Dao container
type Dao struct {
	db *postgres.DB
}

type ResultDaoer interface {
	AddCheckDomainResult(context.Context, structs.CheckResult) (int, error)
}

// NewDAO return Dao instance
func NewDAO(db *postgres.DB) *Dao {
	return &Dao{
		db: db,
	}
}

func (dao *Dao) AddCheckDomainResult(ctx context.Context, r structs.CheckResult) (int, error) {
	log.Debug().
		Str("domain", r.Domain).
		Msg("insert new check result")
	createdOn := time.Now()
	var id int
	err := dao.db.QueryRow(
		ctx,
		`INSERT INTO public."checks" (domain, status_code, error, created_on) VALUES ($1,$2,$3, $4) RETURNING id;`,
		&r.Domain,
		&r.StatusCode,
		&r.Error,
		createdOn,
	).Scan(&id)
	if err != nil {
		log.Error().
			Err(err).
			Str("domain", r.Domain).
			Int("status_code", r.StatusCode).
			Str("error", r.Error).
			Msg("failed to insert new check result")
		return id, err
	}
	return id, nil
}
