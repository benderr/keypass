package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/benderr/keypass/internal/server/domain/record"
	"github.com/benderr/keypass/pkg/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type recordRepository struct {
	db  *sql.DB
	log logger.Logger
}

// New return instance for manipulate record in database
func New(db *sql.DB, log logger.Logger) *recordRepository {
	return &recordRepository{db: db, log: log}
}

func (u *recordRepository) GetByID(ctx context.Context, id string) (*record.Record, error) {

	row := u.db.QueryRowContext(ctx, "SELECT id, user_id, info, data_type, version, updated_at, meta from records WHERE id = $1", id)
	var ord record.Record
	err := row.Scan(&ord.ID, &ord.UserID, &ord.Info, &ord.DataType, &ord.Version, &ord.UpdatedAt, &ord.Meta)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, record.ErrNotFound
		}

		return nil, err
	}

	return &ord, nil
}

func (u *recordRepository) GetByUser(ctx context.Context, userID string) ([]record.Record, error) {
	recordlist := make([]record.Record, 0)

	rows, err := u.db.QueryContext(ctx, "SELECT id, user_id, info, data_type, version, updated_at, meta from records WHERE user_id=$1 ORDER BY updated_at desc", userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var ord record.Record
		err = rows.Scan(&ord.ID, &ord.UserID, &ord.Info, &ord.DataType, &ord.Version, &ord.UpdatedAt, &ord.Meta)
		if err != nil {
			return nil, err
		}

		recordlist = append(recordlist, ord)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return recordlist, nil
}

func (u *recordRepository) Create(ctx context.Context, userID string, info []byte, dataType record.DataType, meta string) (bool, error) {
	_, err := u.db.ExecContext(ctx, `INSERT INTO records (user_id, info, data_type, version, meta) VALUES ($1, $2, $3, 1, $4)`, userID, info, dataType, meta)
	if err != nil {
		var perr *pgconn.PgError
		if errors.As(err, &perr) && perr.Code == pgerrcode.UniqueViolation {
			return false, record.ErrAlreadyExist
		}

		return false, err
	}

	return true, nil
}

func (u *recordRepository) Update(ctx context.Context, ID string, info []byte, meta string) error {
	_, err := u.db.ExecContext(ctx, `UPDATE records SET info=$1, version=records.version + 1, updated_at=NOW(), meta = $2 WHERE id=$3`, info, meta, ID)
	return err
}

func (u *recordRepository) Delete(ctx context.Context, ID string) error {
	_, err := u.db.ExecContext(ctx, `DELETE FROM records WHERE id=$1`, ID)
	return err
}
