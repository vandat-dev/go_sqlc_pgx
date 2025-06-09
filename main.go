package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"go_sqlc_pgx/internal/db"
)

func main() {
	dbURL := "postgresql://admin:123123@172.17.0.1:5002/pgx_go?sslmode=disable"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}
	queries := db.New(pool)

	r := gin.Default()

	// API tạo user
	r.POST("/users", func(c *gin.Context) {
		var req struct {
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Convert phone string to pgtype.Text
		var phone pgtype.Text
		if req.Phone != "" {
			phone = pgtype.Text{String: req.Phone, Valid: true}
		}

		user, err := queries.CreateUser(c, db.CreateUserParams{
			Name:  req.Name,
			Phone: phone,
			Email: req.Email,
		})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, user)
	})

	// API lấy danh sách user
	r.GET("/users", func(c *gin.Context) {
		users, err := queries.ListUsers(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, users)
	})

	// API lấy user theo ID
	r.GET("/users/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID"})
			return
		}

		user, err := queries.GetUserByID(c, int32(id))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, user)
	})

	r.Run(":8080")
}
