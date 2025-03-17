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
	ErrUpdateFail = errors.New("Device update failed")
	ErrDeleteFail = errors.New("Device deletion failed")
)

type DeviceRepository interface {
	Create(Device) (*Device, error)
	Read(id int) (*Device, error)
	Update(Device) (*Device, error)
	Delete(id int) (*Device, error)

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

// Returns either a successfully persisted Device, or a generic error
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
		log.Printf("Error in INSERT-ing %s: %s", d, err)
		return nil, ErrCreateFail
	}

	return &persisted, nil
}

// Returns 1 row, or an ErrNotFound
func (r *DeviceRepositoryPSQL) Read(id int) (*Device, error) {
	retrieved := &Device{}
	err := r.conn.QueryRow(
		r.ctx,
		"SELECT id, rate, model FROM devices WHERE id = $1",
		id,
	).Scan(
		&retrieved.ID,
		&retrieved.Rate,
		&retrieved.Model,
	)

	if err == pgx.ErrNoRows {
		log.Printf("Error in SELECT-ing Device with id %d: %s", id, err)
		return nil, ErrNotFound
	}

	return retrieved, nil
}

// Updates the Device specified by the Device.ID of the argument
// Returns the persisted updated Device, or a generic error
// while the real update error is logged
func (r *DeviceRepositoryPSQL) Update(d Device) (*Device, error) {
	persisted := Device{}
	err := r.conn.QueryRow(
		r.ctx,
		"UPDATE devices SET rate = $2, model = $3 WHERE id = $1 RETURNING id, rate, model",
		d.ID, d.Rate, d.Model,
	).Scan(
		&persisted.ID,
		&persisted.Rate,
		&persisted.Model,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Printf("Error in UPDATE-ing %s: %s", d, err)
		return nil, ErrUpdateFail
	}

	return &persisted, nil
}

func (r *DeviceRepositoryPSQL) Delete(id int) (*Device, error) {
	deleted := Device{}
	err := r.conn.QueryRow(
		r.ctx,
		"DELETE FROM devices WHERE id = $1 RETURNING id, rate, model",
		id,
	).Scan(
		&deleted.ID,
		&deleted.Rate,
		&deleted.Model,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Printf("Error in DELETE-ing Device with id %d: %s", id, err)
		return nil, ErrDeleteFail
	}

	return &deleted, nil
}

func (r *DeviceRepositoryPSQL) Close() {
	err := r.conn.Close(r.ctx)
	if err != nil {
		panic(err)
	}
}
