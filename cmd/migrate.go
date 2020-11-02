package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const MIGRATEION_VERSION_FORMAT = "20060102150405"

type migrationLog struct {
}

func (l *migrationLog) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *migrationLog) Verbose() bool {
	return false
}

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migration commands",
}

var migrateCreateCmd = &cobra.Command{
	Use:     "create migration_name",
	Short:   "Create migration file",
	Args:    cobra.ExactArgs(1),
	Example: "create add_table",
	Run: func(cmd *cobra.Command, args []string) {
		version := migrateGenerateVersion()

		var files = []string{
			fmt.Sprintf("%v/%v_%v.%v.sql", migrateGetDir(), version, args[0], "up"),
			fmt.Sprintf("%v/%v_%v.%v.sql", migrateGetDir(), version, args[0], "down"),
		}

		for _, f := range files {
			if err := migrateCreateFile(f); err != nil {
				fmt.Println(fmt.Sprintf("Create migration file failed: %v ", f), err)
			} else {
				fmt.Println(fmt.Sprintf("Created migration file: %v ", f))
			}
		}
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up [step]",
	Short: "Up action",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := initMigrate()
		if err != nil {
			log.Fatal(err)
		}

		if len(args) > 0 {
			steps := args[0]
			if steps != "" {
				n, err := strconv.ParseUint(steps, 10, 64)
				if err != nil {
					log.Fatal(err)
				}
				if err := m.Steps(int(n)); err != nil {
					if err != migrate.ErrNoChange {
						log.Fatal(err)
					} else {
						log.Println(err)
					}
				}
				return
			}
		} else {
			if err := m.Up(); err != nil {
				if err != migrate.ErrNoChange {
					log.Fatal(err)
				} else {
					log.Println(err)
				}
			}
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down [step]",
	Short: "Down action",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := initMigrate()
		if err != nil {
			log.Fatal(err)
		}

		num := -1

		if len(args) > 0 {
			steps := args[0]
			if steps != "" {
				n, err := strconv.ParseUint(steps, 10, 64)
				if err != nil {
					log.Fatal(err)
				}
				num = -int(n)
			}
		}

		if err := m.Steps(num); err != nil {
			if err != migrate.ErrNoChange {
				log.Fatal(err)
			} else {
				log.Println(err)
			}
		}
	},
}

var migrateForceCmd = &cobra.Command{
	Use:   "force VERSION",
	Args:  cobra.ExactArgs(1),
	Short: "Force version to version",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := initMigrate()
		if err != nil {
			log.Fatal(err)
		}

		v, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatal("error: can't read version argument version")
		}

		if err := m.Force(int(v)); err != nil {
			log.Fatal(err)
		}
	},
}

var migrateVerCmd = &cobra.Command{
	Use:   "version",
	Short: "Current migration version",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := initMigrate()
		if err != nil {
			log.Fatal(err)
		}

		version, dirty, err := m.Version()
		if err != nil {
			if err != migrate.ErrNilVersion {
				log.Fatal(err)
			} else {
				log.Println("No migrate applied")
			}
		}

		if dirty {
			fmt.Printf("%v (dirty)\n", version)
		} else {
			fmt.Println(version)
		}
	},
}

var migrateRefreshCmd = &cobra.Command{
	Use:   "refresh VERSION",
	Short: "Refresh specific version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stderr, "", log.LstdFlags)

		ver, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			logger.Fatal("error: can't read version argument version")
		}

		files, err := filepath.Glob(fmt.Sprintf("%s/%d_*.sql", migrateGetDir(), ver))
		if err != nil {
			logger.Fatalf("find migrate files error: %s", err)
		}

		if files == nil {
			logger.Fatalf("no migrate files found with version %d", ver)
		}

		newVersion := migrateGenerateVersion()
		for _, v := range files {
			newName := strings.Replace(v, fmt.Sprintf("%d", ver), newVersion, 1)
			if err := os.Rename(v, newName); err != nil {
				logger.Fatalf("rename file error: %s", err)
			}

			fmt.Printf("%s ==> %s\n", v, newName)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateCreateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateForceCmd)
	migrateCmd.AddCommand(migrateVerCmd)
	migrateCmd.AddCommand(migrateRefreshCmd)
}

func initMigrate() (*migrate.Migrate, error) {
	m, err := migrate.New(migrateDirToSource(migrateGetDir()), migrateGetDatabaseURL())
	if err != nil {
		return nil, err
	}

	m.Log = &migrationLog{}

	return m, nil
}

func migrateGenerateVersion() string {
	return time.Now().Format(MIGRATEION_VERSION_FORMAT)
}

func migrateDirToSource(dir string) string {
	return "file://" + dir
}

func migrateGetDir() string {
	return strings.TrimSuffix(viper.GetString("migration.path"), "/")
}

func migrateGetDatabaseURL() string {
	return viper.GetString("migration.dsn")
}

func migrateCreateFile(path string) error {
	if _, err := os.Create(path); err != nil {
		return err
	}
	return nil
}
