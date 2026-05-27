package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/FudSy/Diploma/docs"
	"github.com/FudSy/Diploma/internal/pkg/handler"
	"github.com/FudSy/Diploma/internal/pkg/models"
	"github.com/FudSy/Diploma/internal/pkg/repository"
	"github.com/FudSy/Diploma/internal/pkg/repository/postgres"
	"github.com/FudSy/Diploma/internal/pkg/server"
	"github.com/FudSy/Diploma/internal/pkg/service"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// @title Diploma API
// @version 1.0
// @description API for authentication, resources, and bookings.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	if err := initConfig(); err != nil {
		log.Fatal().Msgf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal().Msgf("error loading env variables: %s", err.Error())
	}

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatal().Msgf("failed to initialize db: %s", err.Error())
	}

	db.AutoMigrate(
		&models.User{},
		&models.Resource{},
		&models.Booking{},
		&models.ResourceType{},
		&models.ResourceTypeOption{},
	)

	// Replace AutoMigrate-created unique index on calendar_token with a partial
	// unique index that ignores empty strings (so multiple users may have no token yet).
	db.Exec(`DROP INDEX IF EXISTS idx_users_calendar_token`)
	db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS uniq_users_calendar_token ON users (calendar_token) WHERE calendar_token <> ''`)

	// Ensure FK constraints on bookings cascade on delete (so deleting a resource
	// or user wipes its bookings). AutoMigrate doesn't change existing FK rules.
	db.Exec(`ALTER TABLE bookings DROP CONSTRAINT IF EXISTS fk_resources_bookings`)
	db.Exec(`ALTER TABLE bookings ADD CONSTRAINT fk_resources_bookings FOREIGN KEY (resource_id) REFERENCES resources(id) ON UPDATE CASCADE ON DELETE CASCADE`)
	db.Exec(`ALTER TABLE bookings DROP CONSTRAINT IF EXISTS fk_users_bookings`)
	db.Exec(`ALTER TABLE bookings ADD CONSTRAINT fk_users_bookings FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE`)

	defaultTypes := []string{"MEETING_ROOM", "CAR", "DEVICE"}
	for _, name := range defaultTypes {
		db.Where(models.ResourceType{Name: name}).FirstOrCreate(&models.ResourceType{Name: name})
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Fatal().Msgf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Print("Server Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Server Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error().Msgf("error occured on server shutting down: %s", err.Error())
	}

	sqlDB, err := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Error().Msgf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("internal/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
