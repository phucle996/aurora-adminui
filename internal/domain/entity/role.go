package entity

type Role struct {
	ID              string
	Name            string
	Scope           string
	Description     string
	UserCount       int
	PermissionCount int
}
