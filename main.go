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

// Simple response struct
type ProductResponse struct {
	ID          int32      `json:"ID"`
	Name        string     `json:"Name"`
	Description *string    `json:"Description"`
	Price       *float64   `json:"Price"`
	UserID      int32      `json:"UserID"`
	CreatedAt   *time.Time `json:"CreatedAt"`
	UpdatedAt   *time.Time `json:"UpdatedAt"`
	User        struct {
		UserName  string  `json:"UserName"`
		UserEmail string  `json:"UserEmail"`
		UserPhone *string `json:"UserPhone"`
	} `json:"User"`
}

type UserWithProductsResponse struct {
	ID        int32      `json:"ID"`
	Name      string     `json:"Name"`
	Phone     *string    `json:"Phone"`
	Email     string     `json:"Email"`
	CreatedAt *time.Time `json:"CreatedAt"`
	Products  []struct {
		ID          *int32     `json:"ID"`
		Name        *string    `json:"Name"`
		Description *string    `json:"Description"`
		Price       *float64   `json:"Price"`
		CreatedAt   *time.Time `json:"CreatedAt"`
		UpdatedAt   *time.Time `json:"UpdatedAt"`
	} `json:"Products"`
}

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

	// =============== USER ENDPOINTS ===============

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

	// API lấy danh sách users với products của họ
	r.GET("/users/with-products", func(c *gin.Context) {
		rows, err := queries.ListUsersWithProducts(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		userMap := make(map[int32]*UserWithProductsResponse)

		for _, row := range rows {
			user, exists := userMap[row.UserID]
			if !exists {
				user = &UserWithProductsResponse{
					ID:    row.UserID,
					Name:  row.UserName,
					Email: row.UserEmail,
				}
				if row.UserPhone.Valid {
					user.Phone = &row.UserPhone.String
				}
				if row.UserCreatedAt.Valid {
					user.CreatedAt = &row.UserCreatedAt.Time
				}
				userMap[row.UserID] = user
			}

			// Nếu có product
			if row.ProductID.Valid {
				product := struct {
					ID          *int32     `json:"ID"`
					Name        *string    `json:"Name"`
					Description *string    `json:"Description"`
					Price       *float64   `json:"Price"`
					CreatedAt   *time.Time `json:"CreatedAt"`
					UpdatedAt   *time.Time `json:"UpdatedAt"`
				}{}

				if row.ProductID.Valid {
					product.ID = &row.ProductID.Int32
				}
				if row.ProductName.Valid {
					product.Name = &row.ProductName.String
				}
				if row.ProductDescription.Valid {
					product.Description = &row.ProductDescription.String
				}
				if row.ProductCreatedAt.Valid {
					product.CreatedAt = &row.ProductCreatedAt.Time
				}
				if row.ProductUpdatedAt.Valid {
					product.UpdatedAt = &row.ProductUpdatedAt.Time
				}

				user.Products = append(user.Products, product)
			}
		}

		var result []UserWithProductsResponse
		for _, user := range userMap {
			result = append(result, *user)
		}

		c.JSON(200, result)
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

	// =============== PRODUCT ENDPOINTS ===============

	// API tạo product
	r.POST("/products", func(c *gin.Context) {
		var req struct {
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Price       float64 `json:"price"`
			UserID      int32   `json:"user_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Convert description to pgtype.Text
		var description pgtype.Text
		if req.Description != "" {
			description = pgtype.Text{String: req.Description, Valid: true}
		}

		// Convert price to pgtype.Numeric
		var price pgtype.Numeric
		err := price.Scan(req.Price)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid price format"})
			return
		}
		println(err)

		product, err := queries.CreateProduct(c, db.CreateProductParams{
			Name:        req.Name,
			Description: description,
			Price:       price,
			UserID:      req.UserID,
		})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, product)
	})

	// API lấy danh sách products với user info
	r.GET("/products", func(c *gin.Context) {
		products, err := queries.ListProductsWithUsers(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, products)
	})

	// API lấy product theo ID với user info
	r.GET("/products/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid product ID"})
			return
		}

		p, err := queries.GetProductWithUserByID(c, int32(id))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		result := ProductResponse{
			ID:     p.ID,
			Name:   p.Name,
			UserID: p.UserID,
		}

		if p.Description.Valid {
			result.Description = &p.Description.String
		}
		if p.CreatedAt.Valid {
			result.CreatedAt = &p.CreatedAt.Time
		}
		if p.UpdatedAt.Valid {
			result.UpdatedAt = &p.UpdatedAt.Time
		}

		result.User.UserName = p.UserName
		result.User.UserEmail = p.UserEmail
		if p.UserPhone.Valid {
			result.User.UserPhone = &p.UserPhone.String
		}

		c.JSON(200, result)
	})

	// API lấy products theo user ID
	r.GET("/users/:id/products", func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID"})
			return
		}

		products, err := queries.GetProductsByUserID(c, int32(userID))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, products)
	})

	// API cập nhật product
	r.PUT("/products/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid product ID"})
			return
		}

		var req struct {
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Price       float64 `json:"price"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Convert description to pgtype.Text
		var description pgtype.Text
		if req.Description != "" {
			description = pgtype.Text{String: req.Description, Valid: true}
		}

		// Convert price to pgtype.Numeric
		var price pgtype.Numeric
		err = price.Scan(strconv.FormatFloat(req.Price, 'f', 2, 64))
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid price format"})
			return
		}

		product, err := queries.UpdateProduct(c, db.UpdateProductParams{
			ID:          int32(id),
			Name:        req.Name,
			Description: description,
			Price:       price,
		})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, product)
	})

	// API xóa product
	r.DELETE("/products/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid product ID"})
			return
		}

		err = queries.DeleteProduct(c, int32(id))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Product deleted successfully"})
	})

	r.Run(":8080")
}
