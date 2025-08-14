package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Migration struct {
	Version int
	Name    string
	Content string
}

func Run() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	// Load environment variables
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Printf("Warning: Could not load .env.dev file: %v", err)
	}

	switch command {
	case "up":
		runMigrations()
	case "down":
		if len(os.Args) < 3 {
			log.Fatal("Down command requires number of migrations to rollback")
		}
		steps, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Invalid number of steps:", err)
		}
		rollbackMigrations(steps)
	case "status":
		showMigrationStatus()
	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Create command requires migration name")
		}
		createMigration(os.Args[2])
	case "reset":
		resetDatabase()
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("RecursiveDine Database Migration Tool")
	fmt.Println("=====================================")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/migrate/main.go up                    - Run all pending migrations")
	fmt.Println("  go run cmd/migrate/main.go down <steps>          - Rollback migrations")
	fmt.Println("  go run cmd/migrate/main.go status                - Show migration status")
	fmt.Println("  go run cmd/migrate/main.go create <name>         - Create new migration file")
	fmt.Println("  go run cmd/migrate/main.go reset                 - Reset database (DANGER!)")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go up")
	fmt.Println("  go run cmd/migrate/main.go down 1")
	fmt.Println("  go run cmd/migrate/main.go create add_user_table")
	fmt.Println("  go run cmd/migrate/main.go status")
}

func getDBConnection() *sql.DB {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "recursive_dine")
	sslmode := getEnv("DB_SSL_MODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	return db
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func ensureMigrationTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	if _, err := db.Exec(query); err != nil {
		log.Fatal("Failed to create migration table:", err)
	}
}

func loadMigrations() ([]Migration, error) {
	var migrations []Migration

	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %v", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Extract version from filename (e.g., "001_initial_schema.sql")
		parts := strings.SplitN(file.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("Warning: Invalid migration filename: %s", file.Name())
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join("migrations", file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %v", file.Name(), err)
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    strings.TrimSuffix(file.Name(), ".sql"),
			Content: string(content),
		})
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func getAppliedMigrations(db *sql.DB) (map[int]bool, error) {
	applied := make(map[int]bool)

	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, nil
}

func runMigrations() {
	db := getDBConnection()
	defer db.Close()

	ensureMigrationTable(db)

	migrations, err := loadMigrations()
	if err != nil {
		log.Fatal("Failed to load migrations:", err)
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		log.Fatal("Failed to get applied migrations:", err)
	}

	fmt.Println("Running migrations...")

	count := 0
	for _, migration := range migrations {
		if applied[migration.Version] {
			fmt.Printf("✓ Migration %03d_%s already applied\n", migration.Version, migration.Name[4:])
			continue
		}

		fmt.Printf("→ Applying migration %03d_%s...", migration.Version, migration.Name[4:])

		tx, err := db.Begin()
		if err != nil {
			log.Fatal("Failed to begin transaction:", err)
		}

		// Execute migration
		if _, err := tx.Exec(migration.Content); err != nil {
			tx.Rollback()
			log.Fatalf("Failed to execute migration %d: %v", migration.Version, err)
		}

		// Record migration
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version); err != nil {
			tx.Rollback()
			log.Fatalf("Failed to record migration %d: %v", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalf("Failed to commit migration %d: %v", migration.Version, err)
		}

		fmt.Println(" ✓")
		count++
	}

	if count == 0 {
		fmt.Println("No pending migrations found.")
	} else {
		fmt.Printf("Successfully applied %d migration(s).\n", count)
	}
}

func rollbackMigrations(steps int) {
	db := getDBConnection()
	defer db.Close()

	applied, err := getAppliedMigrations(db)
	if err != nil {
		log.Fatal("Failed to get applied migrations:", err)
	}

	fmt.Printf("Rolling back %d migration(s)...\n", steps)

	// Get applied versions in descending order
	var versions []int
	for version := range applied {
		versions = append(versions, version)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(versions)))

	if len(versions) < steps {
		log.Printf("Warning: Only %d migrations can be rolled back", len(versions))
		steps = len(versions)
	}

	for i := 0; i < steps; i++ {
		version := versions[i]
		fmt.Printf("→ Rolling back migration %03d...", version)

		tx, err := db.Begin()
		if err != nil {
			log.Fatal("Failed to begin transaction:", err)
		}

		// Remove migration record
		if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version); err != nil {
			tx.Rollback()
			log.Fatalf("Failed to remove migration record %d: %v", version, err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalf("Failed to commit rollback %d: %v", version, err)
		}

		fmt.Println(" ✓")
	}

	fmt.Printf("Successfully rolled back %d migration(s).\n", steps)
	fmt.Println("Note: This tool only removes migration records. Manual schema changes may be required.")
}

func showMigrationStatus() {
	db := getDBConnection()
	defer db.Close()

	ensureMigrationTable(db)

	migrations, err := loadMigrations()
	if err != nil {
		log.Fatal("Failed to load migrations:", err)
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		log.Fatal("Failed to get applied migrations:", err)
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	if len(migrations) == 0 {
		fmt.Println("No migrations found.")
		return
	}

	for _, migration := range migrations {
		status := "✗ Pending"
		if applied[migration.Version] {
			status = "✓ Applied"
		}

		fmt.Printf("%s  %03d_%s\n", status, migration.Version, migration.Name[4:])
	}

	pendingCount := 0
	for _, migration := range migrations {
		if !applied[migration.Version] {
			pendingCount++
		}
	}

	fmt.Printf("\nTotal migrations: %d\n", len(migrations))
	fmt.Printf("Applied: %d\n", len(applied))
	fmt.Printf("Pending: %d\n", pendingCount)
}

func createMigration(name string) {
	// Get next version number
	migrations, err := loadMigrations()
	if err != nil {
		log.Fatal("Failed to load migrations:", err)
	}

	nextVersion := 1
	if len(migrations) > 0 {
		lastMigration := migrations[len(migrations)-1]
		nextVersion = lastMigration.Version + 1
	}

	// Create filename
	filename := fmt.Sprintf("%03d_%s.sql", nextVersion, name)
	filepath := filepath.Join("migrations", filename)

	// Create migration template
	template := fmt.Sprintf(`-- Migration: %s
-- Created: %s

-- Add your migration SQL here
-- Example:
-- CREATE TABLE example (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- Don't forget to add appropriate indexes:
-- CREATE INDEX idx_example_name ON example(name);
`, name, time.Now().Format("2006-01-02 15:04:05"))

	// Write file
	if err := ioutil.WriteFile(filepath, []byte(template), 0644); err != nil {
		log.Fatal("Failed to create migration file:", err)
	}

	fmt.Printf("Created migration file: %s\n", filename)
	fmt.Println("Edit the file to add your migration SQL, then run:")
	fmt.Println("  go run cmd/migrate/main.go up")
}

func resetDatabase() {
	fmt.Print("Are you sure you want to reset the database? This will delete ALL data! (yes/no): ")
	var response string
	fmt.Scanln(&response)

	if response != "yes" {
		fmt.Println("Database reset cancelled.")
		return
	}

	db := getDBConnection()
	defer db.Close()

	fmt.Println("Resetting database...")

	// Drop all tables
	tables := []string{
		"payments", "order_items", "orders", "menu_items",
		"menu_categories", "tables", "users", "schema_migrations",
	}

	for _, table := range tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: Failed to drop table %s: %v", table, err)
		} else {
			fmt.Printf("✓ Dropped table: %s\n", table)
		}
	}

	fmt.Println("Database reset completed. Run 'go run cmd/migrate/main.go up' to recreate tables.")
}
