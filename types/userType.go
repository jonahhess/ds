package types

type UserType int

const (
	Guest UserType = iota
	loggedInUser
	Admin
)