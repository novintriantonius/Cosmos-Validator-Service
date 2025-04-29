package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed migrations/*.up.sql
var migrationsFS embed.FS

// RunMigrations executes all SQL migration files in the migrations directory
func RunMigrations(db *sql.DB) error {
	// Read all files from the embedded filesystem
	files, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}
	
	// Filter and sort migration files
	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)
	
	// Execute each migration in order
	for _, filename := range migrationFiles {
		// Read the migration file
		content, err := fs.ReadFile(migrationsFS, filepath.Join("migrations", filename))
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %v", filename, err)
		}
		
		// Split the content into individual statements
		statements := strings.Split(string(content), ";")
		
		// Execute each statement
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			
			// Execute the statement
			_, err := db.Exec(stmt)
			if err != nil {
				return fmt.Errorf("error executing migration %s: %v", filename, err)
			}
		}
		
		fmt.Printf("Successfully applied migration: %s\n", filename)
	}
	
	return nil
} 