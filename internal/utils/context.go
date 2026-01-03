package utils

type contextKey string

var (
	TxKey       contextKey = "tx"
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)
