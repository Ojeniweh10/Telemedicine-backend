package database

import (
	"context"
	"fmt"
	"log"

	"telemed/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	host     = config.Db().Host
	user     = config.Db().User
	password = config.Db().Password
	dbname   = config.Db().Name
)

func NewConnection() *pgxpool.Pool {
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, host, dbname)
	log.Println("connecting to database at", databaseUrl)
	ctx := context.Background()

	poolConfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		panic(err)
	}

	poolConfig.MaxConns = 5
	dbPool, err := pgxpool.ConnectConfig(ctx, poolConfig)

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")
	return dbPool
}

func Insert(db *pgxpool.Pool, table string, data map[string]any) error {
	var columns, placeholders string
	var values []any
	i := 1
	for k, v := range data {
		if i > 1 {
			columns += ","
			placeholders += ","
		}
		columns += k
		placeholders += fmt.Sprintf("$%d", i)
		values = append(values, v)
		i++
	}
	sql := fmt.Sprintf("INSERT INTO %v (%s) VALUES (%s)", table, columns, placeholders)
	_, err := db.Exec(context.Background(), sql, values...)
	if err != nil {
		return err
	}

	return nil
}

func Update(db *pgxpool.Pool, table string, data map[string]any, condition map[string]any) error {
	var setClause string
	var values []any
	i := 1
	for k, v := range data {
		if i > 1 {
			setClause += ", "
		}
		setClause += fmt.Sprintf("%s = $%d", k, i)
		values = append(values, v)
		i++
	}
	whereClause := ""
	for k, v := range condition {
		if whereClause == "" {
			whereClause = " WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += fmt.Sprintf("%s = $%d", k, i)
		values = append(values, v)
		i++
	}
	sql := fmt.Sprintf("UPDATE %s SET %s%s", table, setClause, whereClause)
	_, err := db.Exec(context.Background(), sql, values...)
	if err != nil {
		return err
	}

	return nil
}
