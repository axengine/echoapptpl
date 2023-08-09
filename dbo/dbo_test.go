package dbo

import "testing"

func TestSync(t *testing.T) {
	db := New("user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
	db.Sync()
}
