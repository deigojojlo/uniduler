package utils
import (
    "context"
    "github.com/jackc/pgx/v5"
    "github.com/joho/godotenv"
    "os"
)

func Connect() (*pgx.Conn, error) {
    err := godotenv.Load()
	if err != nil {
        panic(err)
    }

	// Récupère les variables d'environnement
	dbUser := os.Getenv("LOGIN")
	dbPassword := os.Getenv("PASSWD")
	dbName := os.Getenv("DBNAME")
    conn, err := pgx.Connect(context.Background(), "postgres://" + dbUser + ":" + dbPassword + "@127.0.0.1:5432/"+ dbName)
    if err != nil {
        return nil, err
    }
    return conn, nil
}

func AddData(conn *pgx.Conn,event *Event){
	sqlStatement :=
	`INSERT INTO events
	(name,groups,summary,location,startDate,endDate,year, dayOfTheWeek, type, parcours) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
    ctx := context.Background()
	_, err := conn.Exec(ctx,sqlStatement, event.Name, event.Groups, event.Summary, event.Location, event.Start, event.End, event.Year, event.DayOfTheWeek,event.Type,event.Parcours)
	if err != nil {
		panic(err)
	}
}

