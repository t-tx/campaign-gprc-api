package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	scriptPath = "./scripts"
)

func Get(host, dbName, username, passord string) *pgxpool.Pool {
	return initPosgresDB(host, dbName, username, passord)
}

func initPosgresDB(host, dbName, username, passord string) *pgxpool.Pool {
	ctx := context.Background()

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s", username, passord, host, dbName)

	newDB, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database.")
	}
	return newDB
}

func RunMigrations(dbPath string, db *pgxpool.Pool) {
	files, err := os.ReadDir(scriptPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read directory")
	}

	var sqlScipts []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			sqlScipts = append(sqlScipts, file.Name())
		}
	}
	sort.Slice(sqlScipts, func(i, j int) bool {
		return sqlScipts[i] < sqlScipts[j]
	})

	lastScript, err := getLastMigrateScriptFile(db)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get last migrate script file")
	}
	for _, filename := range sqlScipts {
		if lastScript < filename {
			scriptContent, err := os.ReadFile(fmt.Sprintf("%s/%s", scriptPath, filename))
			if err != nil {
				log.Fatal().Err(err).Msg("failed to read file")
			}
			tx, err := db.Begin(bgCtx)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to begin transaction")
			}
			_, err = tx.Exec(bgCtx, string(scriptContent))
			if err != nil {
				log.Fatal().Err(err).Msg("failed to execute script")
			}
			var mv = &migrateValue{
				Content: string(scriptContent),
				File:    filename,
			}
			jsonValue, err := json.Marshal(mv)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to marshal json")
			}

			sqlStatement := `
				INSERT INTO global_config (key, value, created_at, scope)
				VALUES
					('migrate', $1, $2, 'system')
				ON CONFLICT (key, scope) DO UPDATE SET value = $1, created_at = $2
			`

			_, err = tx.Exec(bgCtx, sqlStatement, jsonValue, GetTimeNowString())
			if err != nil {
				log.Fatal().Err(err).Str("filename", filename).Msg("failed to insert global config")
			}
			err = tx.Commit(bgCtx)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to commit transaction")
			}

			log.Info().Str("file", filename).Msg("Executed")
		}
	}
}

var bgCtx = context.Background()

type migrateValue struct {
	File    string `json:"file"`
	Content string `json:"content"`
}

func getLastMigrateScriptFile(db *pgxpool.Pool) (string, error) {
	_, err := db.Exec(bgCtx, `CREATE TABLE IF NOT EXISTS global_config (
    id SERIAL PRIMARY KEY,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at TEXT NOT NULL,
    scope TEXT NOT NULL,
    CONSTRAINT global_config_key_scope_key UNIQUE (key, scope)
);`)
	if err != nil {
		return "", err
	}
	var value string
	err = db.QueryRow(bgCtx, `SELECT value FROM global_config WHERE key = $1 and scope = $2 limit 1`, "migrate", "system").Scan(&value)
	if err != nil {
		if err.Error() == "sql: no rows in result set" || err.Error() == "no rows in result set" {
			return "", nil
		}
		return "", err
	}

	var mv migrateValue
	err = json.Unmarshal([]byte(value), &mv)
	if err != nil {
		return "", err
	}
	return mv.File, nil
}
func GetTimeNowString() string {
	return time.Now().UTC().Format(time.RFC3339)
}
func GetTimeString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
func ParseTime(timeString string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeString)
}
