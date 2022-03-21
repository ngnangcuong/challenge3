package models

type Role struct {
	Name 	string `json:"name"`
	Permission 	string `json:"permission"`
}

type RoleRepo interface {
	Create(role Role) (error)
	Find(name string) (Role, error)
}