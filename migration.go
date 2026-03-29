package gormigrate

import (
	"fmt"
	"sort"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Migration struct {
	Version     int
	Description string
	Up          func(*gorm.DB) error
	Down        func(*gorm.DB) error
}

type Registry struct {
	migrations []Migration
}

func NewRegistry() *Registry {
	return &Registry{
		migrations: make([]Migration, 0),
	}
}

func (r *Registry) Register(m Migration) {
	r.migrations = append(r.migrations, m)
}

func (r *Registry) GetMigrations() []Migration {
	sorted := make([]Migration, len(r.migrations))
	copy(sorted, r.migrations)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Version < sorted[j].Version
	})

	return sorted
}

func (r *Registry) Run(db *gorm.DB, targetVersion int) error {
	db.Logger = logger.Default.LogMode(logger.Info)

	if err := r.validate(targetVersion); err != nil {
		return fmt.Errorf("migration validation error: %w", err)
	}

	fmt.Printf("Expected schema version: %d\n", targetVersion)
	fmt.Println("Applying migrations...")

	migrations := r.GetMigrations()
	appliedCount := 0

	for _, m := range migrations {
		if m.Version <= targetVersion {
			fmt.Printf("Applying migration %d: %s\n", m.Version, m.Description)

			if err := apply(db, m); err != nil {
				return fmt.Errorf("error applying migration %d: %w", m.Version, err)
			}

			appliedCount++
		}
	}

	if appliedCount == 0 {
		fmt.Println("No migrations to apply.")
	} else {
		fmt.Printf("%d migration(s) applied successfully!\n", appliedCount)
	}

	fmt.Printf("Schema is now at version %d\n", targetVersion)
	return nil
}

func (r *Registry) Rollback(db *gorm.DB, currentVersion, targetVersion int) error {
	db.Logger = logger.Default.LogMode(logger.Info)

	if targetVersion < 0 {
		return fmt.Errorf("invalid target version: %d", targetVersion)
	}

	if targetVersion >= currentVersion {
		return fmt.Errorf("target version (%d) must be less than current version (%d)", targetVersion, currentVersion)
	}

	migrations := r.GetMigrations()

	for i := len(migrations) - 1; i >= 0; i-- {
		m := migrations[i]
		if m.Version > targetVersion && m.Version <= currentVersion {
			if m.Down == nil {
				return fmt.Errorf("migration %d does not have a Down function", m.Version)
			}

			fmt.Printf("Reverting migration %d: %s\n", m.Version, m.Description)

			tx := db.Begin()
			if err := m.Down(tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("error reverting migration %d: %w", m.Version, err)
			}

			if err := tx.Commit().Error; err != nil {
				return fmt.Errorf("error committing transaction: %w", err)
			}
		}
	}

	fmt.Printf("Schema rolled back to version %d\n", targetVersion)
	return nil
}

func apply(db *gorm.DB, m Migration) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := m.Up(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *Registry) validate(targetVersion int) error {
	migrations := r.GetMigrations()

	if len(migrations) == 0 && targetVersion > 0 {
		return fmt.Errorf("target version is %d but no migrations were registered", targetVersion)
	}

	for v := 1; v <= targetVersion; v++ {
		found := false
		for _, m := range migrations {
			if m.Version == v {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("migration version %d was not registered (expected target version: %d)", v, targetVersion)
		}
	}

	for _, m := range migrations {
		if m.Version > targetVersion {
			return fmt.Errorf("migration version %d was registered but target version is %d", m.Version, targetVersion)
		}
	}

	return nil
}

var defaultRegistry = NewRegistry()

func Register(m Migration) {
	defaultRegistry.Register(m)
}

func Run(db *gorm.DB, targetVersion int) error {
	return defaultRegistry.Run(db, targetVersion)
}

func Rollback(db *gorm.DB, currentVersion, targetVersion int) error {
	return defaultRegistry.Rollback(db, currentVersion, targetVersion)
}
