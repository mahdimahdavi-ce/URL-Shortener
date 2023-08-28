package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testStoreService = &StoreService{}

func init() {
	testStoreService = InitializeStore()
}

func TestStoreInit(t *testing.T) {
	assert.True(t, testStoreService.postgreSqlClient != nil)
	assert.True(t, testStoreService.redisClient != nil)
}

func TestInsertionAndRetrieval(t *testing.T) {
	originalUrl := "https://digikala.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortURL := "Jsz4k57oAX"

	SaveUrlMapping(shortURL, originalUrl)

	retrivedOriginalUrl, _ := RetriveOriginalUrlFromDb(shortURL)

	assert.Equal(t, originalUrl, retrivedOriginalUrl)

	retrivedOriginalUrl, _ = RetrieveOriginalUrl(shortURL)

	assert.Equal(t, originalUrl, retrivedOriginalUrl)

	retrivedOriginalUrl, _ = testStoreService.redisClient.Get(context.Background(), shortURL).Result()

	assert.Equal(t, originalUrl, retrivedOriginalUrl)

}
