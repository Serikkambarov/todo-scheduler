package db

import (
    "database/sql"
    "fmt"
    "os"

    _ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX idx_date ON scheduler(date);
`

func Init(dbFile string) error {
    install := false
    if _, err := os.Stat(dbFile); os.IsNotExist(err) {
        install = true
    }

    db, err := sql.Open("sqlite", dbFile)
    if err != nil {
        return fmt.Errorf("cannot open DB: %v", err)
    }

    if install {
        _, err = db.Exec(schema)
        if err != nil {
            db.Close()
            return fmt.Errorf("cannot create schema: %v", err)
        }
    }

    DB = db
    return nil
}
