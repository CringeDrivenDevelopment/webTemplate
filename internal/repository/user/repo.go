package userRepo

import "github.com/jackc/pgx/v5/pgxpool"

type UserRepository struct {
	pgxpool *pgxpool.Pool
}

func New(pgxpool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pgxpool: pgxpool,
	}
}
