package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ama struct {
	AskedBy    string    `json:"askedBy"`
	Question   string    `json:"question"`
	Answer     string    `json:"answer"`
	AnsweredAt time.Time `json:"answeredAt"`
}

var amas = []ama{
	{AskedBy: "johnson", Question: "Daily perfume ?", Answer: "Molecule02", AnsweredAt: time.Now()},
}

func main() {
	router := gin.Default()
	router.GET("/api/questions", getAmas)
	router.POST("/api/questions", postAmas)

	router.Run("localhost:8080")
}

func getAmas(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, amas)
}

func postAmas(c *gin.Context) {
	var newAma ama
	if err := c.BindJSON(&newAma); err != nil {
		return
	}

	amas = append(amas, newAma)
	c.IndentedJSON(http.StatusCreated, newAma)
}
