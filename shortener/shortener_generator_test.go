package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortUrlGenerator(t *testing.T) {
	initialLink_1 := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.htmle0dba740-fc4b-4977-872c-d360239e6b1a"
	initialLink_2 := "https://www.eddywm.com/lets-build-a-url-shortener-in-go-with-redis-part-2-storage-layer/e0dba740-fc4b-4977-872c-d360239e6b1a"
	initialLink_3 := "https://spectrum.ieee.org/automaton/robotics/home-robots/hello-robots-stretch-mobile-manipulatore0dba740-fc4b-4977-872c-d360239e6b1a"

	shortLink_1 := GenerateShortLink(initialLink_1)
	shortLink_2 := GenerateShortLink(initialLink_2)
	shortLink_3 := GenerateShortLink(initialLink_3)

	assert.Equal(t, shortLink_1, "jTa4L57P")
	assert.Equal(t, shortLink_2, "d66yfx7N")
	assert.Equal(t, shortLink_3, "dhZTayYQ")
}
