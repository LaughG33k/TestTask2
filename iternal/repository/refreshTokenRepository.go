package repository

import (
	"context"
	"database/sql"

	"gopkg.in/reform.v1"
)

type RefreshToken struct {
	Db *reform.DB
}

func (r *RefreshToken) CreateRefreshToken(ctx context.Context, oldRefreshToken, refreshToken, ownerUuid string, tokenTimeLife int64) error {

	if oldRefreshToken != "" {

		tx, err := r.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})

		if err != nil {
			return err
		}

		defer tx.Rollback()

		if _, err := tx.ExecContext(ctx, "delete from refresh_tokens where token=$1 and owner_uuid=$2;", oldRefreshToken, ownerUuid); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, "insert into refresh_tokens(token, owner_uuid, time_end_of_life) values($1, $2, $3);", refreshToken, ownerUuid, tokenTimeLife); err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		return nil
	}

	if _, err := r.Db.ExecContext(ctx, "insert into refresh_tokens(token, owner_uuid, time_end_of_life) values($1, $2, $3);", refreshToken, ownerUuid, tokenTimeLife); err != nil {
		return err
	}

	return nil
}

/*
this func finds required refresh token and return uuid of owner and token timelive
*/
func (r *RefreshToken) FindRefreshToken(ctx context.Context, refreshToken string) (string, int64, error) {

	ownerUuid := ""
	var timeLife int64 = 0

	if err := r.Db.QueryRowContext(ctx, "select owner_uuid, time_end_of_life from refresh_tokens where token=$1;", refreshToken).Scan(&ownerUuid, &timeLife); err != nil {
		return "", 0, err
	}

	return ownerUuid, timeLife, nil

}
