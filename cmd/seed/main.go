package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/FudSy/Diploma/internal/pkg/models"
	"github.com/FudSy/Diploma/internal/pkg/repository/postgres"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatal().Msgf("config error: %s", err.Error())
	}
	_ = godotenv.Load()

	host := viper.GetString("db.host")
	if h := os.Getenv("DB_HOST"); h != "" {
		host = h
	}
	port := viper.GetString("db.port")
	if p := os.Getenv("DB_PORT"); p != "" {
		port = p
	}

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     host,
		Port:     port,
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatal().Msgf("db connect error: %s", err.Error())
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Resource{},
		&models.Booking{},
		&models.ResourceType{},
		&models.ResourceTypeOption{},
	); err != nil {
		log.Fatal().Msgf("migrate error: %s", err.Error())
	}

	db.Exec(`DROP INDEX IF EXISTS idx_users_calendar_token`)
	db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS uniq_users_calendar_token ON users (calendar_token) WHERE calendar_token <> ''`)

	// ── Resource types ──
	for _, name := range []string{"MEETING_ROOM", "CAR", "DEVICE"} {
		db.Where(models.ResourceType{Name: name}).FirstOrCreate(&models.ResourceType{Name: name})
	}

	// ── Users ──
	users := []struct {
		Login, Email, Name, Surname, Password, Role string
	}{
		{"admin", "admin@corp.local", "Анна", "Администратор", "admin123", "ADMIN"},
		{"ivan", "ivan@corp.local", "Иван", "Петров", "user1234", "USER"},
		{"maria", "maria@corp.local", "Мария", "Сидорова", "user1234", "USER"},
		{"alex", "alex@corp.local", "Алексей", "Иванов", "user1234", "USER"},
	}

	userMap := map[string]models.User{}
	for _, u := range users {
		var existing models.User
		if err := db.Where("login = ?", u.Login).First(&existing).Error; err == nil {
			userMap[u.Login] = existing
			continue
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal().Msgf("bcrypt: %s", err.Error())
		}
		user := models.User{
			Login:        u.Login,
			Email:        u.Email,
			Name:         u.Name,
			Surname:      u.Surname,
			FullName:     strings.TrimSpace(u.Name + " " + u.Surname),
			PasswordHash: string(hash),
			Role:         u.Role,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Fatal().Msgf("create user %s: %s", u.Login, err.Error())
		}
		userMap[u.Login] = user
		fmt.Printf("user created: %s (%s) — password: %s\n", u.Login, u.Role, u.Password)
	}

	// ── Resources ──
	resources := []models.Resource{
		{Name: "Переговорная «Орбита»", Description: "Большой зал, проектор, видеосвязь", Type: "MEETING_ROOM", Capacity: 12, IsActive: true, Location: "Этаж 3, к. 301"},
		{Name: "Переговорная «Меркурий»", Description: "Маленькая комната для 1-on-1", Type: "MEETING_ROOM", Capacity: 4, IsActive: true, Location: "Этаж 3, к. 305"},
		{Name: "Переговорная «Сатурн»", Description: "Средняя комната с TV", Type: "MEETING_ROOM", Capacity: 8, IsActive: true, Location: "Этаж 4, к. 412"},
		{Name: "Toyota Camry", Description: "Корпоративный автомобиль 2022", Type: "CAR", Capacity: 5, IsActive: true, Location: "Парковка-1, место 7"},
		{Name: "Hyundai Tucson", Description: "Внедорожник для дальних поездок", Type: "CAR", Capacity: 5, IsActive: true, Location: "Парковка-1, место 8"},
		{Name: "MacBook Pro 16\"", Description: "M3 Max, 64 GB / 1 TB", Type: "DEVICE", Capacity: 1, IsActive: true, Location: "IT-склад, шкаф A"},
		{Name: "Проектор Epson EB-W42", Description: "1080p, HDMI/USB-C", Type: "DEVICE", Capacity: 1, IsActive: true, Location: "IT-склад, шкаф B"},
		{Name: "VR-комплект Meta Quest 3", Description: "Для презентаций и демо", Type: "DEVICE", Capacity: 1, IsActive: false, Location: "IT-склад, на ремонте"},
	}

	resMap := map[string]models.Resource{}
	for _, r := range resources {
		var existing models.Resource
		if err := db.Where("name = ?", r.Name).First(&existing).Error; err == nil {
			resMap[r.Name] = existing
			continue
		}
		if err := db.Create(&r).Error; err != nil {
			log.Fatal().Msgf("create resource %s: %s", r.Name, err.Error())
		}
		resMap[r.Name] = r
		fmt.Printf("resource created: %s [%s]\n", r.Name, r.Type)
	}

	// ── Bookings ──
	// Time anchor — round current time to nearest hour for predictable demo.
	now := time.Now().Truncate(time.Hour)

	type bk struct {
		resourceName string
		userLogin    string
		startOffset  time.Duration
		duration     time.Duration
		status       string
	}

	bookings := []bk{
		// Past bookings (yesterday and earlier)
		{"Переговорная «Орбита»", "ivan", -26 * time.Hour, 1 * time.Hour, "CONFIRMED"},
		{"Toyota Camry", "maria", -48 * time.Hour, 3 * time.Hour, "CONFIRMED"},
		{"Переговорная «Меркурий»", "alex", -72 * time.Hour, 1 * time.Hour, "CANCELLED"},
		{"Переговорная «Сатурн»", "ivan", -5 * 24 * time.Hour, 2 * time.Hour, "CONFIRMED"},
		{"MacBook Pro 16\"", "maria", -10 * 24 * time.Hour, 8 * time.Hour, "CONFIRMED"},

		// Live (started a bit ago, still running)
		{"Переговорная «Орбита»", "maria", -30 * time.Minute, 90 * time.Minute, "CONFIRMED"},
		{"Toyota Camry", "alex", -15 * time.Minute, 4 * time.Hour, "CONFIRMED"},

		// Upcoming
		{"Переговорная «Меркурий»", "ivan", 2 * time.Hour, 1 * time.Hour, "CONFIRMED"},
		{"Переговорная «Сатурн»", "alex", 4 * time.Hour, 90 * time.Minute, "CONFIRMED"},
		{"Hyundai Tucson", "ivan", 24 * time.Hour, 8 * time.Hour, "CONFIRMED"},
		{"MacBook Pro 16\"", "ivan", 48 * time.Hour, 9 * time.Hour, "CONFIRMED"},
		{"Проектор Epson EB-W42", "maria", 26 * time.Hour, 2 * time.Hour, "CONFIRMED"},
		{"Переговорная «Орбита»", "alex", 72 * time.Hour, 2 * time.Hour, "CONFIRMED"},
		{"Переговорная «Сатурн»", "ivan", 96 * time.Hour, 1 * time.Hour, "CANCELLED"},
		{"Toyota Camry", "ivan", 7 * 24 * time.Hour, 4 * time.Hour, "CONFIRMED"},
	}

	// Clear ephemeral bookings only if explicitly asked
	if os.Getenv("SEED_RESET_BOOKINGS") == "1" {
		if err := db.Exec(`DELETE FROM bookings`).Error; err != nil {
			log.Fatal().Msgf("reset bookings: %s", err.Error())
		}
		fmt.Println("bookings reset")
	}

	for _, b := range bookings {
		resource, ok := resMap[b.resourceName]
		if !ok {
			log.Warn().Msgf("skip booking: resource %q not found", b.resourceName)
			continue
		}
		user, ok := userMap[b.userLogin]
		if !ok {
			log.Warn().Msgf("skip booking: user %q not found", b.userLogin)
			continue
		}

		start := now.Add(b.startOffset)
		end := start.Add(b.duration)

		// Idempotent: skip if a booking already exists for this user/resource/start.
		var count int64
		db.Model(&models.Booking{}).
			Where("user_id = ? AND resource_id = ? AND start_time = ?", user.ID, resource.ID, start).
			Count(&count)
		if count > 0 {
			continue
		}

		booking := models.Booking{
			UserID:     user.ID,
			ResourceID: resource.ID,
			StartTime:  start,
			EndTime:    end,
			Status:     b.status,
		}
		if err := db.Create(&booking).Error; err != nil {
			log.Warn().Msgf("create booking failed: %s", err.Error())
			continue
		}
		fmt.Printf("booking: %s / %s — %s [%s]\n",
			user.Login, resource.Name, start.Format("2006-01-02 15:04"), b.status)
	}

	fmt.Println("\nSeed complete.")
}

func initConfig() error {
	viper.AddConfigPath("internal/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
