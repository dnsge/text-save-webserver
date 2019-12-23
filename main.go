package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const charset = "ABCDEFGHJK1234567890"
const charLimit = 5000
const timeToLive = time.Minute * 30

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RedisDatabase() gin.HandlerFunc {
	redisAddr := os.Getenv("REDIS_ADDRESS")
	redisPass := os.Getenv("REDIS_PASS")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       redisDB,
	})

	return func(c *gin.Context) {
		c.Set("redis", client)
		c.Next()
	}
}

func SaveRender(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{})
}

func Save(c *gin.Context) {
	client := c.MustGet("redis").(*redis.Client)

	text := c.PostForm("text")
	if text == "" {
		c.Redirect(http.StatusFound, "/")
	} else if len(text) > charLimit {
		c.HTML(http.StatusForbidden, "too_long.html", gin.H{
			"charLimit": charLimit,
		})
	} else {
		var code string
		var key string

		for {
			code = RandomStringWithCharset(6, charset)
			key = "textSave:texts:" + code
			if n, err := client.Exists(key).Result(); err != nil {
				panic(err)
			} else if n == 0 { // 0 means key doesn't exist, so we can use specific code
				break
			}
		}

		err := client.Set(key, text, timeToLive).Err()
		if err != nil {
			panic(err)
		}

		c.Redirect(http.StatusFound, "/"+code)
	}
}

func Get(c *gin.Context) {
	client := c.MustGet("redis").(*redis.Client)

	code := c.Param("code")
	key := "textSave:texts:" + code
	value, err := client.Get(key).Result()
	if err == redis.Nil {
		c.HTML(http.StatusNotFound, "not_found.html", gin.H{
			"code": code,
		})
	} else if err != nil {
		panic(err)
	} else {
		ttl, err := client.TTL(key).Result()
		if err != nil {
			panic(err)
		}

		minutes := int(ttl.Minutes())
		seconds := int(ttl.Seconds()) % 60
		expiresString := fmt.Sprintf("%d min %d sec", minutes, seconds)

		// c.Header("Cache-Control", "public, max-age=" + strconv.Itoa(int(ttl.Seconds())))
		c.HTML(http.StatusOK, "record.html", gin.H{
			"text":    value,
			"code":    code,
			"expires": expiresString,
		})
	}
}

func main() {
	r := gin.Default()
	r.Use(RedisDatabase())

	r.LoadHTMLGlob("templates/*")

	r.GET("/", SaveRender)
	r.POST("/", Save)
	r.GET("/:code", Get)
	r.Run()
}
