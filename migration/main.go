package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
  _ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
	"online_chess/util"
	"strconv"
)

func main() {
  args := os.Args
	var (
		migrate_version string
		confirm_migrate string
		input_hint string
	)

	if len(args) == 1 {
		input_hint = "Database will migrate to latest version, are you sure?(y/n)"
	} else {
		migrate_version = args[1]
		input_hint = fmt.Sprintf("Database will migrate to %s version, are you sure?(y/n)", migrate_version)
	}
	for confirm_migrate != "n" && confirm_migrate != "y" {
		fmt.Printf(input_hint)
		fmt.Scanf("%s", &confirm_migrate)
	}
	if confirm_migrate == "n" {
		fmt.Println("Cancel migration !!!")
		return
	}

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file", err)
		return
	}
	db_user, user_ok := util.ReadEnvVariable("DB_USER")
	db_pwd, pwd_ok := util.ReadEnvVariable("DB_PWD")
	if !user_ok || !pwd_ok {
		fmt.Println("db env is empty!!!")
		return
	}
	m, err := migrate.New(
		"file://migration",
		"mysql://" + db_user + ":" + db_pwd + "@tcp(localhost:3306)/online_war_chess?query")
	if err != nil {
		fmt.Println("migrate.New fail error: ", err)
		return
	}

	if migrate_version == "" {
		fmt.Println("directly migrate up to latest version")
		if err := m.Up(); err != nil {
			fmt.Println("migration execute error", err)
		} else {
			fmt.Println("migration success")
		}
	} else {
		version, _, _ := m.Version()
		fmt.Println("currently active migration version: ", version)
		step, _ := strconv.Atoi(migrate_version)
		step = step - int(version)
		if err := m.Steps(step); err != nil {
			fmt.Println("migration execute error", err)
		} else {
			fmt.Println("migration success")
		}
	}
}