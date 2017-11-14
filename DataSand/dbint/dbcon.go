package dbint

import ("database/sql"
        _ "github.com/lib/pq"
        "log"
        "fmt"
)

const (
        host     = "localhost"
        port     = 5432
        user     = "postgres"
        password = "cisco123"
        dbname   = "messages"
)

type DBData struct {
        *sql.DB
}


func (dbc *DBData) Connect(){
        dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
                host, port, user, password, dbname)
        db, err := sql.Open("postgres", dbinfo)
        dbc.DB = db
        CheckError(err)
        _, e:= db.Exec("CREATE TABLE if NOT EXISTS packets (source char(36)," +
                "dest char(36), data bytea)")
        CheckError(e)
}

func CheckError(err error){
        if err!=nil{
                log.Fatal("Error: ", err)
        }
}
