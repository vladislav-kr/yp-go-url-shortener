// dbkeeper отвечает за хранилище в postgres
package dbkeeper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/lib/cryptoutils"
)

// Ошибки уровня хранилища
var (
	ErrAlreadyExists = errors.New("the value already exists")
	ErrURLRemoved    = errors.New("url has already been deleted")
)

// DBKeeper хранит подключения к БД.
type DBKeeper struct {
	db     *sql.DB
	dbPool *pgxpool.Pool
	log    *zap.Logger
}

// NewDBKeeper конструктор DBKeeper.
func NewDBKeeper(log *zap.Logger, db *sql.DB, dbPool *pgxpool.Pool) *DBKeeper {
	return &DBKeeper{
		db:     db,
		dbPool: dbPool,
		log:    log,
	}
}

// PostURL сохранение сокращенного URL.
func (k *DBKeeper) PostURL(ctx context.Context, url string, userID string) (string, error) {

	id, err := cryptoutils.GenerateRandomString(10)
	if err != nil {
		return "", err
	}

	sqlStatement := `
		INSERT INTO shortened_url (short_url, original_url, user_id)
		VALUES ($1, $2, $3)`

	_, err = k.db.ExecContext(
		ctx,
		sqlStatement,
		id, url, NullUserID(userID),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			sqlStatement := `SELECT short_url FROM shortened_url WHERE original_url=$1;`
			row := k.db.QueryRowContext(
				ctx,
				sqlStatement,
				url,
			)
			id := ""
			if err := row.Scan(&id); err != nil {
				return "", err
			}

			return id, ErrAlreadyExists
		}

		return "", err
	}
	return id, nil
}

// SaveURLS массовое сохранение URL.
func (k *DBKeeper) SaveURLS(ctx context.Context, urls []models.BatchRequest, userID string) ([]models.BatchResponse, error) {
	tx, err := k.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			k.log.Error("fail rollback",
				zap.Error(err),
			)
		}
	}()

	sqlStatement := `
		INSERT INTO shortened_url(short_url, original_url, user_id)
		VALUES ($1, $2, $3)`

	stmt, err := tx.PrepareContext(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			k.log.Error("failed close of prepared statement",
				zap.Error(err),
			)
		}
	}()

	batchResp := make([]models.BatchResponse, 0, len(urls))

	for _, url := range urls {
		id, err := cryptoutils.GenerateRandomString(10)
		if err != nil {
			return nil, err
		}

		_, err = stmt.ExecContext(ctx, id, url.OriginalURL, NullUserID(userID))

		if err != nil {
			return nil, err
		}

		batchResp = append(batchResp, models.BatchResponse{
			CorrelationID: url.CorrelationID,
			ShortURL:      id,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return batchResp, nil
}

// GetURL чтение оригинального URL.
func (k *DBKeeper) GetURL(ctx context.Context, id string) (string, error) {
	sqlStatement := `SELECT original_url, is_deleted FROM shortened_url WHERE short_url=$1;`

	row := k.db.QueryRowContext(ctx, sqlStatement, id)

	var fullURL string
	var deleted bool

	err := row.Scan(&fullURL, &deleted)
	if err != nil {
		return "", fmt.Errorf("records for the key %s do not exist", id)
	}
	if deleted {
		return "", ErrURLRemoved
	}

	return fullURL, nil
}

// GetURLS список сокращенных URL пользователя.
func (k *DBKeeper) GetURLS(ctx context.Context, userID string) ([]models.MassURL, error) {
	sqlStatement := `
		SELECT
			short_url,
			original_url
		FROM
			shortened_url
		WHERE
			user_id = $1;`

	rows, err := k.db.QueryContext(ctx, sqlStatement, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []models.MassURL{}, nil

		default:
			return nil, err
		}
	}

	if rows.Err() != nil {
		return nil, err
	}

	urls := []models.MassURL{}
	for rows.Next() {
		url := models.MassURL{}
		err = rows.Scan(&url.ShortURL, &url.OriginalURL)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

// NullUserID создает sql.NullString
func NullUserID(userID string) sql.NullString {
	var valid bool
	if len(userID) > 0 {
		valid = true
	}

	return sql.NullString{
		String: userID,
		Valid:  valid,
	}

}

// DeleteURLS удаление URL
func (k *DBKeeper) DeleteURLS(ctx context.Context, shortURLS []models.DeleteURL) {

	query := `
		UPDATE shortened_url
		SET
			is_deleted = true
		WHERE
			short_url = @shortURL
			AND user_id = @userID`

	batch := &pgx.Batch{}
	for _, url := range shortURLS {
		args := pgx.NamedArgs{
			"shortURL": url.ShortURL,
			"userID":   url.UserID,
		}
		batch.Queue(query, args)
	}

	results := k.dbPool.SendBatch(ctx, batch)
	defer results.Close()

	for _, url := range shortURLS {
		_, err := results.Exec()
		if err != nil {
			k.log.Error("failed to delete url",
				zap.String("url", url.ShortURL),
				zap.Error(err),
			)
		}
	}

	if err := results.Close(); err != nil {
		k.log.Error("failed to close response batch",
			zap.Error(err),
		)
	}
}
