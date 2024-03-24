package main

import (
	"chatbot/clients/telegram"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	TOKEN = "bot1:key1234567890"
)

type Buf struct {
	Reviews []telegram.Review
	idx     int
	m       sync.Mutex
}

var buf Buf

func (b *Buf) next() (review *telegram.Review) {
	b.m.Lock()
	defer b.m.Unlock()

	review = &b.Reviews[b.idx]
	review.ReviewId = uuid.NewString()
	review.CreateTs = time.Now().Unix()

	b.idx++
	if b.idx == len(b.Reviews) {
		b.idx = 0
	}
	return
}

func main() {
	rBytes, err := os.ReadFile("reviews.json")
	if err != nil {
		log.Fatalln("read reviews.json failed:", err)
		return
	}
	if err = json.Unmarshal(rBytes, &buf.Reviews); err != nil {
		log.Fatalln("parse reviews.json failed:", err)
		return
	}

	s := gin.Default()
	s.GET("/"+TOKEN, func(c *gin.Context) {
		c.JSON(http.StatusOK, buf.next())
	})
	s.POST("/"+TOKEN, func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		log.Println("recv ", string(body))
		c.JSON(http.StatusOK, nil)
	})
	s.Run(":10000")
}
