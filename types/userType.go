package types

type UserType int

const (
	Guest UserType = iota
	LoggedInUser
)
