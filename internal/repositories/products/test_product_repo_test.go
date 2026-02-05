package productrepositories

import (
	"fmt"
	"testing"
)

type productTest struct {
	ID    int
	Stock int
}

func TestString(t *testing.T) {
	product := []productTest{
		{ID: 1, Stock: 300},
		{ID: 5, Stock: 400},
	}
	sql := "UPDATE products SET stock = CASE id "
	args := []interface{}{}
	ids := []interface{}{}

	for _, p := range product {
		sql += "WHEN ? THEN ? "
		args = append(args, p.ID, p.Stock)
		ids = append(ids, p.ID)
	}

	sql += "END WHERE id IN (?)"
	args = append(args, ids)
	fmt.Println(args)
	fmt.Println(sql)
	fmt.Println(args...)
}
