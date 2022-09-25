package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func DatabaseConnect() {

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	databaseUrl := "postgres://postgres:siburian@localhost:5000/personal-web"

	var err error
	Conn, err = pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		// fmt.Println("Koneksi gagal", err)
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)

	}

	fmt.Println("Database connected")
}
