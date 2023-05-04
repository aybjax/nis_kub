package app_db

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) {
	version1(db)
}

func version1(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS students(
		"id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
		"name" VARCHAR(256),
		"first_name" VARCHAR(256),
		"last_name" VARCHAR(256),
		"course_ids" VARCHAR(256) ARRAY NOT NULL DEFAULT '{}',
		CONSTRAINT u_constraint UNIQUE (name, first_name, last_name)
	)`)

	if err != nil {
		panic(fmt.Sprint(err))
	}
}
