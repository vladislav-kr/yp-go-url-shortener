package dbkeeper

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/lib/cryptoutils"
)

type DBKeeper struct {
	db *sql.DB
}

func NewDBKeeper(db *sql.DB) *DBKeeper {
	return &DBKeeper{
		db: db,
	}
}

func (k *DBKeeper) PostURL(ctx context.Context, url string) (string, error) {

	id, err := cryptoutils.GenerateRandomString(10)
	if err != nil {
		return "", err
	}

	sqlStatement := `
		INSERT INTO shortened_url (short_url, original_url)
		VALUES ($1, $2)`

	_, err = k.db.ExecContext(
		ctx,
		sqlStatement,
		id, url,
	)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (k *DBKeeper) GetURL(ctx context.Context, id string) (string, error) {
	sqlStatement := `SELECT original_url FROM shortened_url WHERE short_url=$1;`

	row := k.db.QueryRowContext(ctx, sqlStatement, id)

	var fullURL string

	err := row.Scan(&fullURL)
	if err != nil {
		return "", fmt.Errorf("records for the key %s do not exist", id)
	}
	return fullURL, nil
}
