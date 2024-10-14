package controller

import (
	"context"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hmuir28/go-thepapucoin/models"
)

var validate = validator.New()

func SendMoney() gin.HandlerFunc {
	return func (c *gin.Context) {

		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var transaction models.Transaction

		if err := c.BindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(transaction)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		c.JSON(http.StatusCreated, "Successfully transfered!")
	}
}
