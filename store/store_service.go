package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type StoreService struct {
	redisClient      *redis.Client
	postgreSqlClient *gorm.DB
}

var (
	storeService = new(StoreService)
	ctx          = context.Background()
)

type UrlMapping struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	shortUrl    string
	originalUrl string
}

func loadEnvVariables() {
	err := godotenv.Load("../local.env")
	if err != nil {
		log.Fatal("Error while loading .env file")
	}
}

func loadPostgresUrl() string {
	host := os.Getenv("host")
	dbUsername := os.Getenv("dbUsername")
	dbPassword := os.Getenv("dbPassword")
	dbName := os.Getenv("dbName")
	dbPort := os.Getenv("dbPort")

	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, dbUsername, dbPassword, dbName, dbPort)

	return url
}

func loadRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     os.Getenv("redisAdd"),
		Password: "",
		DB:       0,
	}
}

func InitializeStore() *StoreService {
	loadEnvVariables()
	postgreUrl := loadPostgresUrl()
	redisOptions := loadRedisOptions()

	postgreClient, postgresErr := gorm.Open(postgres.Open(postgreUrl), &gorm.Config{})
	if postgresErr != nil {
		log.Fatal("Failed to connect to the PostgreSQL Database")
	}
	// log
	fmt.Println("PostgreSQL started successfully")

	redisClient := redis.NewClient(redisOptions)
	_, redisErr := redisClient.Ping(ctx).Result()
	if redisErr != nil {
		log.Fatal("Failed to connect to the Redis")
	}
	// log
	fmt.Println("Redis started successfully")

	storeService.postgreSqlClient = postgreClient
	storeService.redisClient = redisClient

	return storeService
}

func saveUrlMapping(shortUrl string, originalUrl string) error {
	newShortUrlRecord := UrlMapping{
		shortUrl:    shortUrl,
		originalUrl: originalUrl,
	}

	result := storeService.postgreSqlClient.Create(newShortUrlRecord)
	if result.Error != nil {
		fmt.Printf("Failed saving the urls: OriginalUrl: %v  shortUrl: %v", originalUrl, shortUrl)
		return errors.New(fmt.Sprintf("Failed saving the urls: OriginalUrl: %v  shortUrl: %v", originalUrl, shortUrl))
	}

	return nil
}

func RetrieveInitialUrl(shortUrl string) (string, error) {
	originalUrl, err := storeService.redisClient.Get(ctx, shortUrl).Result()
	if err == redis.Nil {
		var urlMappingRecord UrlMapping
		// fetching from database
		result := storeService.postgreSqlClient.Where("shortUrl = ?", shortUrl).First(urlMappingRecord)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return "", errors.New("The specified url dose not exist in database")
			} else {
				return "", errors.New(fmt.Sprintf("Sth went wrong while reading from database: %v", result.Error))
			}
		}

		originalUrl := urlMappingRecord.originalUrl
		// inserting into redis
		redisErr := storeService.redisClient.Set(ctx, shortUrl, originalUrl, 0).Err()
		if redisErr != nil {
			fmt.Printf("Failed saving urls in redis: shortUrl: %v, originalUrl: %v\n", shortUrl, originalUrl)
		}

		return originalUrl, nil

	} else {
		return originalUrl, nil
	}
}
