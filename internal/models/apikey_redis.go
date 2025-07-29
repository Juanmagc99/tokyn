package models

type APIKeyRedis struct {
	ID      string `redis:"primaryKey"`
	KeyHash string `redis:"keyHash"`
	Name    string `redis:"name"`
}
