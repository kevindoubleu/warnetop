package device

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
)

// DB-agnostic errors
var (
	ErrNotFound   = errors.New("Device not found")
	ErrCreateFail = errors.New("Device creation failed")
)

type DeviceRepository interface {
	Create(Device) (*Device, error)
	Read(int) (*Device, error)

	Close()
}

type DeviceRepositoryPSQL struct {
	conn *pgx.Conn
	ctx  context.Context
}

func NewDeviceRepositoryPSQL() *DeviceRepositoryPSQL {
	conn_str := "user=warnetop_admin password=secret host=localhost port=5432 dbname=warnetop"
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, conn_str)
	if err != nil {
		panic(err)
	}

	return &DeviceRepositoryPSQL{
		conn: conn,
		ctx:  ctx,
	}
}

// Returns either a successfully persisted Device or a generic error
// while the real creation error is logged
func (r *DeviceRepositoryPSQL) Create(d Device) (*Device, error) {
	persisted := Device{}
	err := r.conn.QueryRow(
		r.ctx,
		"INSERT INTO devices(rate, model) VALUES($1, $2) RETURNING id, rate, model",
		d.Rate, d.Model,
	).Scan(
		&persisted.ID,
		&persisted.Rate,
		&persisted.Model,
	)

	if err != nil {
		log.Printf("Error in inserting %s: %s", d, err)
		return nil, ErrCreateFail
	}

	return &persisted, nil
}

func (r *DeviceRepositoryPSQL) Read(id int) (*Device, error) {
	row := r.conn.QueryRow(r.ctx, "SELECT id, rate, model FROM devices WHERE id = $1", id)
	retrieved := &Device{}
	err := row.Scan(&retrieved.ID, &retrieved.Rate, &retrieved.Model)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	return retrieved, nil
}

func (r *DeviceRepositoryPSQL) Close() {
	err := r.conn.Close(r.ctx)
	if err != nil {
		panic(err)
	}
}
