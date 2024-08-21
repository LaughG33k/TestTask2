package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gopkg.in/reform.v1"
)

type User struct {
	Db *reform.DB
}

func (r *User) Register(ctx context.Context, login, password, permission string) error {

	if _, err := r.Db.ExecContext(ctx, "insert into users(login, password, permission) values($1, crypt($2, gen_salt('md5')), $3);", login, password, permission); err != nil {

		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			return fmt.Errorf(pgError.Code)
		}

		return err

	}

	return nil

}

/*
return uuid and permission of finded user
*/
func (r *User) CheckByLP(ctx context.Context, login, password string) (bool, string, string, error) {

	t := false
	var uuid string
	var perm string

	if err := r.Db.QueryRowContext(ctx, "select uuid, permission,(password=crypt($2, password)) as t from users where login=$1;", login, password).Scan(&uuid, &perm, &t); err != nil {

		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			return t, "", "", fmt.Errorf(pgError.Code)
		}

		return t, "", "", err

	}

	return t, uuid, perm, nil

}

func (r *User) GetAllFields(ctx context.Context, uuid string) (string, string, string, string, error) {

	var login string
	var password string
	var perm string

	if err := r.Db.QueryRowContext(ctx, "select uuid,login, password, permission uuid from users where uuid=$1;", uuid).Scan(&uuid, &login, &password, &perm); err != nil {

		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			return "", "", "", "", fmt.Errorf(pgError.Code)
		}

		return "", "", "", "", err

	}

	return uuid, login, password, perm, nil

}
