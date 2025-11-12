package postgres

import (
	"context"

	"orderService/models"
	"orderService/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// urlExample := "postgres://username:password@localhost:5432/database_name"
type orderRepo struct {
	db *gorm.DB
}

// func New(cfg config.PostgresConfig) *Store {
// 	var err error
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
// 	defer cancel()
// 	sqlDB, err := sql.Open("pgx", cfg.ConnString())
// 	if err != nil {
// 		log.Fatalf("Что-то пошло не так при подкючении к БД %v", err)
// 	}

// 	if err = sqlDB.PingContext(ctx); err != nil {
// 		log.Fatalf("Не удаётся установить соеденинение к БД %v", err)
// 	}
// 	log.Println("internal/event/config/ConfigHandler():Соединение выполнено успешно")

// 	gormDB, err := gorm.Open(postgres.New(postgres.Config{
// 		Conn: sqlDB,
// 	}), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Не удаётся преобразовать sql config в gorm config %v", err)
// 	}
// 	return &Store{db: gormDB}
// }

func (ms *orderRepo) Create(ctx context.Context, order models.Order) (string, error) {
	cl := logger.LoggerFromContext(ctx)
	result := gorm.WithResult()
	err := gorm.G[models.Order](ms.db, result).Create(ctx, &order)
	if err != nil {
		cl.Error(ctx, "Ошибка при создании заказа в БД ", zap.Error(err))
		return "", err
	}
	return order.ID, nil
}

func (ms *orderRepo) GetByID(ctx context.Context, id string) (models.Order, error) {
	cl := logger.LoggerFromContext(ctx)

	// Using numeric primary key
	order, err := gorm.G[models.Order](ms.db).Where("id = ?", id).First(ctx)
	if err != nil {
		cl.Error(ctx, "Ошибка при создании заказа в БД ", zap.Error(err))
		return models.Order{}, err
	}
	return order, nil
}

func (ms *orderRepo) Update(ctx context.Context, order models.Order) error {
	cl := logger.LoggerFromContext(ctx)
	_, err := gorm.G[models.Order](ms.db).Where("id = ?", order.ID).Updates(ctx,
		models.Order{
			Item:     order.Item,
			Quantity: order.Quantity,
		})
	if err != nil {
		cl.Error(ctx, "Ошибка при создании заказа в БД ", zap.Error(err))
		return err
	}
	return err
}

func (ms *orderRepo) Delete(ctx context.Context, id string) error {
	cl := logger.LoggerFromContext(ctx)
	_, err := gorm.G[models.Order](ms.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		cl.Error(ctx, "Ошибка при создании заказа в БД ", zap.Error(err))
		return err
	}
	return err
}

func (ms *orderRepo) List(ctx context.Context) ([]models.Order, error) {
	cl := logger.LoggerFromContext(ctx)
	orders := make([]models.Order, 0, 1)
	result := ms.db.Find(&orders)
	if result.Error != nil {
		cl.Error(ctx, "Ошибка при создании заказа в БД ", zap.Error(result.Error))
		return []models.Order{}, result.Error
	}
	return orders, nil
}
