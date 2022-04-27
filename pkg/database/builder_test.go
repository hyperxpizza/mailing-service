package database

import (
	"fmt"
	"testing"
)

func TestQueryBuilder(t *testing.T) {
	q, allowedVars := buildGetRecipientsQuery("id", 5, 0)
	fmt.Println(q)
	fmt.Println(allowedVars[0], allowedVars[1], allowedVars[2])
	fmt.Println()

	q2, _ := buildGetRecipientsWhereGroupQuery("", "CUSTOMERS", 0, 0)
	fmt.Println(q2)
}
