// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Product struct {
	ID          int32
	Name        string
	Description pgtype.Text
	Price       pgtype.Numeric
	UserID      int32
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type User struct {
	ID        int32
	Name      string
	Phone     pgtype.Text
	Email     string
	CreatedAt pgtype.Timestamptz
}
