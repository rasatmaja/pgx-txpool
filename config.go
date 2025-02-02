package pgxtxpool

import (
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type config struct {
	dsn   url.URL
	query url.Values
}

func (c *config) SetQuery(key, value string) {
	if c.query == nil {
		c.query = url.Values{}
	}
	c.query.Set(key, value)
}

// ParseToPGXConfig will parse dsn and query URL to pgxpool config
func (c *config) ParseToPGXConfig() *pgxpool.Config {
	c.dsn.Scheme = "postgres"
	c.dsn.RawQuery = c.query.Encode()
	config, err := pgxpool.ParseConfig(c.dsn.String())
	if err != nil {
		panic(err)
	}

	return config
}

// Option is a function that can be used to configure the pgxpool
type Option func(*config)

// SetHost will set postgres host and port
func SetHost(host, port string) Option {
	return func(c *config) {
		c.dsn.Host = fmt.Sprintf("%s:%s", host, port)
	}
}

// SetCredential will set postgres user and password
func SetCredential(user, password string) Option {
	return func(c *config) {
		c.dsn.User = url.UserPassword(user, password)
	}
}

// SetDatabase will set postgres database
func SetDatabase(database string) Option {
	return func(c *config) {
		c.dsn.Path = database
	}
}

// WithSSLMode will set sslmode
func WithSSLMode(mode string) Option {
	return func(c *config) {
		c.SetQuery("sslmode", mode)
	}
}

// WithMaxConns will set max connections
func WithMaxConns(maxConns int) Option {
	return func(c *config) {
		c.SetQuery("pool_max_conns", fmt.Sprintf("%d", maxConns))
	}
}

// WithMaxIdleConns will set max idle connections
// ex: "30s", "5m"
func WithMaxIdleConns(maxIdleConns string) Option {
	return func(c *config) {
		c.SetQuery("pool_max_conn_idle_time", maxIdleConns)
	}
}

// WithMaxConnLifetime will set max connection lifetime
// ex: "30s", "5m"
func WithMaxConnLifetime(maxConnLifetime string) Option {
	return func(c *config) {
		c.SetQuery("pool_max_conn_lifetime", maxConnLifetime)
	}
}
