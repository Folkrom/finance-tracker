package model

// AdminStats holds platform-level statistics. Not a DB model.
type AdminStats struct {
	Users            int `json:"users"`
	CategoriesGlobal int `json:"categories_global"`
	CategoriesUser   int `json:"categories_user"`
	Profiles         int `json:"profiles"`
}
