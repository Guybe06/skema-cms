package conduit

const (
	DriverPostgres = "postgres"
	DriverMySQL    = "mysql"

	ErrUnsupportedDriver  = "pilote de base de données non supporté : %s"
	ErrConnectionFailed   = "échec de la connexion à la base de données"
	ErrConnectionTimeout  = "délai de connexion dépassé"
	ErrPingFailed         = "la base de données ne répond pas"

	PoolMaxConns     = 10
	PoolMinConns     = 2
	ConnTTLMinutes   = 30
)
