package store


type Store struct {
	UnimplementedStore
	db *pgx.Conn
}




