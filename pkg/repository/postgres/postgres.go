package postgres

import (
	"context"
	"database/sql"
	"log"
	"orderService/internal/config"
	"orderService/pkg/repository"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository.Repository {
	return &repo{
		db: db,
	}
}

func (s *repo) GetGorm() *gorm.DB {
	return s.db
}

func (r *repo) Order() repository.OrderRepository {
	return &orderRepo{db: r.db}
}
func (r *repo) User() repository.UserRepository {
	return &userRepo{db: r.db}
}

func New(cfg *config.PostgresConfig) *gorm.DB {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	sqlDB, err := sql.Open("pgx", cfg.ConnString())
	if err != nil {
		log.Fatalf("Что-то пошло не так при подкючении к БД %v", err)
	}

	if err = sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("Не удаётся установить соеденинение к БД %v", err)
	}
	log.Println("internal/event/config/ConfigHandler():Соединение выполнено успешно")

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удаётся преобразовать sql config в gorm config %v", err)
	}
	return gormDB
}
