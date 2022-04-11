package database

import "fmt"

const (
	GroupNotFoundErrorMsg     = "group with name: %s was not found in the database"
	RecipientNotFoundErrorMsg = "recipient with email: %s was not found in the database"
)

type NotFoundError struct {
	Name string
	Msg  string
}

func (e *NotFoundError) Error() string {
	return e.Msg
}

func NewGroupNotFoundError(groupName string) *NotFoundError {
	return &NotFoundError{
		Name: groupName,
		Msg:  fmt.Sprintf(GroupNotFoundErrorMsg, groupName),
	}
}

func NewRecipientNotFoundError(email string) *NotFoundError {
	return &NotFoundError{
		Name: email,
		Msg:  fmt.Sprintf(RecipientNotFoundErrorMsg, email),
	}
}
