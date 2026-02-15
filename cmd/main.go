package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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

func main() {

	if err := initConfig(); err != nil {
		log.Fatal().Msgf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
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
	)


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