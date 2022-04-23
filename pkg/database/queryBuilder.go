package database

import "fmt"

const (
	orderBase  = " order by $%d"
	limitBase  = " limit $%d"
	offsetBase = " offset $%d"
	whereGroup = "where g.groupName=$%d"
)

func buildGetRecipientsQuery(order string, limit, offset int64) (string, []interface{}) {

	fragment, allowedVars := getFragmentAndAllowedVars("", order, limit, offset)
	query := fmt.Sprintf("%s %s", getRecipientsBase, fragment)
	return query, allowedVars
}

func buildGetRecipientsWhereGroupQuery(order, group string, limit, offset int64) (string, []interface{}) {
	var base string
	fragment, allowedVars := getFragmentAndAllowedVars(group, order, limit, offset)
	if group == "" {
		base = getRecipientsBase
	} else {
		base = getRecipientsBaseGroupName
	}

	query := fmt.Sprintf("%s %s", base, fragment)
	return query, allowedVars
}

func getFragmentAndAllowedVars(group, order string, limit, offset int64) (string, []interface{}) {
	fragment := ""
	counter := 0
	var allowedVars []interface{}
	if group != "" {
		counter++
		groupString := fmt.Sprintf(whereGroup, counter)
		fragment += groupString
		allowedVars = append(allowedVars, group)
	}

	if order != "" {
		counter++
		orderString := fmt.Sprintf(orderBase, counter)
		fragment += orderString
		allowedVars = append(allowedVars, order)

	}

	if limit > 0 {
		counter++
		limiString := fmt.Sprintf(limitBase, counter)
		fragment += limiString
		allowedVars = append(allowedVars, limit)
	}

	if offset > 0 {
		counter++
		offsetString := fmt.Sprintf(offsetBase, counter)
		fragment += offsetString
		allowedVars = append(allowedVars, offsetString)
	}

	return fragment, allowedVars
}

func buildSearchQuery(phrase string, limit, offset int64) (string, []interface{}) {

}
