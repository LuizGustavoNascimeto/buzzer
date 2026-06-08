// 1. GORM abre a conexão
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// 2. pega o *sql.DB por baixo do GORM
sqlDB, err := db.DB()

// 3. Goose usa essa mesma conexão pra rodar as migrations
goose.Up(sqlDB, "internal/activities/migrations")

// 4. a partir daqui o GORM usa o banco já migrado normalmente
repo := repository.NewGormActivityRepository(db)