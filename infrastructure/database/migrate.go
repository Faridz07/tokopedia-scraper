package database

import (
	"fmt"
	"os"
	"path/filepath"
)

func MigrateDatabase(db Database, path string) {
	err := ensureMigrationTableExists(db)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Failed to read migration folder: %v\n", err)
		return
	}

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			fmt.Printf("Failed to get info for file %s: %v\n", file.Name(), err)
			continue
		}

		if info.IsDir() || filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		status, err := migrationAlreadyExecuted(db, file.Name())
		if err != nil {
			continue
		}
		if status {
			fmt.Printf("Skipping already executed migration %s\n", file.Name())
			continue
		}

		sqlFilePath := filepath.Join(path, file.Name())

		tx, err := db.Begin()
		if err != nil {
			fmt.Printf("Failed to begin transaction for %s: %v\n", file.Name(), err)
			continue
		}

		sqlContent, err := os.ReadFile(sqlFilePath)
		if err != nil {
			fmt.Printf("Failed to read SQL file %s: %v\n", file.Name(), err)
			tx.Rollback()
			continue
		}

		if _, err = tx.Exec(string(sqlContent)); err != nil {
			fmt.Printf("Failed to execute migration %s: %v\n", file.Name(), err)
			tx.Rollback()
			continue
		}

		err = recordMigrationExecution(db, file.Name())
		if err != nil {
			continue
		}

		if err = tx.Commit(); err != nil {
			fmt.Printf("Failed to commit migration %s: %v\n", file.Name(), err)
			continue
		}

		fmt.Printf("Successfully executed migration %s\n", file.Name())
	}
}

func ensureMigrationTableExists(db Database) error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS migrations (
        id SERIAL PRIMARY KEY,
        file_name TEXT UNIQUE NOT NULL,
        executed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating migrations table: %w", err)
	}
	return nil
}

func migrationAlreadyExecuted(db Database, fileName string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM migrations WHERE file_name = $1;`
	err := db.QueryRow(query, fileName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("Error checking migration status for %s: %w\n", fileName, err)
	}
	return count > 0, nil
}

func recordMigrationExecution(db Database, fileName string) error {
	insertSQL := `INSERT INTO migrations (file_name) VALUES ($1);`
	_, err := db.Exec(insertSQL, fileName)
	if err != nil {
		return fmt.Errorf("error recording migration execution: %w", err)
	}
	return nil
}
