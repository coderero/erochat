package mysql

import (
	"database/sql"
	"errors"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type ConnectionPool struct {
	// Database connection pool.
	pool *sql.DB

	// Maximum number of connections in the pool.
	maxConnections int

	// Number of connections currently in the pool.
	conns int

	// Mutex to protect the pool.
	mu *sync.Mutex
}

var (
	// ErrMaxConnections is returned when the maximum number of connections in the pool is reached.
	ErrMaxConnections = errors.New("maximum number of connections in the pool is reached")
)

func NewConnectionPool(dsn string, maxConnections int) (*ConnectionPool, error) {
	// Create a new connection pool.
	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = pool.Ping()
	if err != nil {
		return nil, err
	}

	// Set the maximum number of connections in the pool.
	pool.SetMaxOpenConns(maxConnections)

	return &ConnectionPool{
		pool:           pool,
		maxConnections: maxConnections,
		conns:          0,
		mu:             &sync.Mutex{},
	}, nil
}

func (p *ConnectionPool) Get() (*sql.DB, error) {
	// Lock the pool.
	p.mu.Lock()
	defer p.mu.Unlock()

	// If the number of connections in the pool is less than the maximum number of connections, create a new connection.
	if p.conns < p.maxConnections {
		p.conns++
		return p.pool, nil
	}

	// If the number of connections in the pool is equal to the maximum number of connections, return an error.
	return nil, ErrMaxConnections
}

func (p *ConnectionPool) Release(conn *sql.DB) {
	// Lock the pool.
	p.mu.Lock()
	defer p.mu.Unlock()

	// Decrement the number of connections in the pool.
	p.conns--
}
