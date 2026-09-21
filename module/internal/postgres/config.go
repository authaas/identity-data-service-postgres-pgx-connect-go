//revive:disable:package-comments
package postgres

// Configuration is where the database is reached: the connection string.
type Configuration struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
}
