package pgpoolstorage

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresPoolStorage is an implementation of the Pool interface
// that uses a postgres database to store the data
type PostgresPoolStorage struct {
	db *pgxpool.Pool
}

// NewPostgresPoolStorage creates and initializes an instance of PostgresPoolStorage
func NewPostgresPoolStorage(cfg db.Config) (*PostgresPoolStorage, error) {
	poolDB, err := db.NewSQLDB(cfg)
	if err != nil {
		return nil, err
	}

	return &PostgresPoolStorage{
		db: poolDB,
	}, nil
}

// AddTx adds a transaction to the pool table with the provided state
func (p *PostgresPoolStorage) AddTx(ctx context.Context, tx types.Transaction, state pool.TxState) error {
	hash := tx.Hash().Hex()

	b, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	encoded := hex.EncodeToHex(b)

	b, err = tx.MarshalJSON()
	if err != nil {
		return err
	}
	decoded := string(b)

	receivedAt := time.Now()
	gasPrice := tx.GasPrice().Uint64()
	nonce := tx.Nonce()
	sql := "INSERT INTO pool.txs (hash, encoded, decoded, state, gas_price, nonce, received_at) VALUES($1, $2, $3, $4, $5, $6, $7)"
	if _, err := p.db.Exec(ctx, sql, hash, encoded, decoded, state, gasPrice, nonce, receivedAt); err != nil {
		return err
	}
	return nil
}

// GetTxsByState returns an array of transactions filtered by state
// limit parameter is used to limit amount txs from the db,
// if limit = 0, then there is no limit
func (p *PostgresPoolStorage) GetTxsByState(ctx context.Context, state pool.TxState, limit uint64) ([]pool.Transaction, error) {
	var (
		rows pgx.Rows
		err  error
		sql  string
	)
	if limit == 0 {
		sql = "SELECT encoded, state, received_at FROM pool.txs WHERE state = $1"
		rows, err = p.db.Query(ctx, sql, state.String())
	} else {
		sql = "SELECT encoded, state, received_at FROM pool.txs WHERE state = $1 LIMIT $2"
		rows, err = p.db.Query(ctx, sql, state.String(), limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]pool.Transaction, 0, len(rows.RawValues()))
	for rows.Next() {
		var (
			encoded, state string
			receivedAt     time.Time
		)

		if err := rows.Scan(&encoded, &state, &receivedAt); err != nil {
			return nil, err
		}

		tx := new(pool.Transaction)

		b, err := hex.DecodeHex(encoded)
		if err != nil {
			return nil, err
		}

		if err := tx.UnmarshalBinary(b); err != nil {
			return nil, err
		}

		tx.State = pool.TxState(state)
		tx.ReceivedAt = receivedAt
		txs = append(txs, *tx)
	}

	return txs, nil
}

// CountTransactionsByState get number of transactions
// accordingly to the provided state
func (p *PostgresPoolStorage) CountTransactionsByState(ctx context.Context, state pool.TxState) (uint64, error) {
	sql := "SELECT COUNT(*) FROM pool.txs WHERE state = $1"
	var counter uint64
	err := p.db.QueryRow(ctx, sql, state.String()).Scan(&counter)
	if err != nil {
		return 0, err
	}
	return counter, nil
}

// UpdateTxState updates a transaction state accordingly to the
// provided state and hash
func (p *PostgresPoolStorage) UpdateTxState(ctx context.Context, hash common.Hash, newState pool.TxState) error {
	sql := "UPDATE pool.txs SET state = $1 WHERE hash = $2"
	if _, err := p.db.Exec(ctx, sql, newState, hash.Hex()); err != nil {
		return err
	}
	return nil
}

// UpdateTxsState updates transactions state accordingly to the provided state and hashes
func (p *PostgresPoolStorage) UpdateTxsState(ctx context.Context, hashes []common.Hash, newState pool.TxState) error {
	hh := make([]string, 0, len(hashes))
	for _, h := range hashes {
		hh = append(hh, h.Hex())
	}

	sql := "UPDATE pool.txs SET state = $1 WHERE hash = ANY ($2)"
	if _, err := p.db.Exec(ctx, sql, newState, hh); err != nil {
		return err
	}
	return nil
}

// SetGasPrice allows an external component to define the gas price
func (p *PostgresPoolStorage) SetGasPrice(ctx context.Context, gasPrice uint64) error {
	sql := "INSERT INTO pool.gas_price (price, timestamp) VALUES ($1, $2)"
	if _, err := p.db.Exec(ctx, sql, gasPrice, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

// GetGasPrice returns the current gas price
func (p *PostgresPoolStorage) GetGasPrice(ctx context.Context) (uint64, error) {
	sql := "SELECT price FROM pool.gas_price ORDER BY item_id DESC LIMIT 1"
	rows, err := p.db.Query(ctx, sql)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}

	defer rows.Close()

	gasPrice := uint64(0)

	for rows.Next() {
		err := rows.Scan(&gasPrice)
		if err != nil {
			return 0, err
		}
	}

	return gasPrice, nil
}