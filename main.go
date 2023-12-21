package main

import (
	"flag"

	"main/database"
	"main/database/migration"
	"main/database/seeder"
	"main/internal/factory"
	"main/internal/http"
	"main/internal/middleware"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	database.GetConnection()
	database.GetConnectionMongoDB()
}

func main() {
	// logger
	lumberJack := lumberjack.Logger{
		Filename:   "logs/lumberjack.log",
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	multi := zerolog.MultiLevelWriter(os.Stdout, &lumberJack)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	database.CreateConnection()

	var migrate string
	var seed string

	flag.StringVar(
		&migrate,
		"migrate",
		"none",
		`this argument for check if user want to migrate table, rollback table, or status migration
to use this flag:
	use -migrate=migrate for migrate table
	use -migrate=rollback for rollback table
	use -migrate=status for get status migration`,
	)

	flag.StringVar(
		&seed,
		"seed",
		"none",
		`this argument for check if user want to seed table
to use this flag:
	use -seed=all to seed all table`,
	)

	flag.Parse()

	if migrate == "migrate" {
		migration.Migrate()
	} else if migrate == "rollback" {
		migration.Rollback()
	} else if migrate == "status" {
		migration.Status()
	} else {
		log.Print("No Key Migrate")
	}

	if seed == "all" {
		seeder.NewSeeder().DeleteAll()
		seeder.NewSeeder().SeedAll()
	}

	f := factory.NewFactory()
	e := echo.New()

	middleware.LogMiddlewares(e)

	http.NewHttp(e, f)

	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}
