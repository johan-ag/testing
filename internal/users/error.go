package users

import (
	"errors"
)

var (
	ErrorFindLastInsertedID = errors.New("error to find  the last inserted id")
	ErrorSavingToDB         = errors.New("error saving to db")
)
