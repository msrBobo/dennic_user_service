package entity

import "time"

type User struct {
	Id           string
	UserOrder    uint64
	FirstName    string
	LastName     string
	BirthDate    string
	PhoneNumber  string
	Password     string
	Gender       string
	RefreshToken string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Admin struct {
	Id            string
	AdminOrder    int64
	Role          string
	FirstName     string
	LastName      string
	BirthDate     string
	PhoneNumber   string
	Email         string
	Password      string
	Gender        string
	Salary        float32
	Biography     string
	StartWorkYear string
	EndWorkYear   string
	WorkYears     uint64
	RefreshToken  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CheckFieldReq struct {
	Value string
	Field string
}

type CheckFieldResp struct {
	Status bool
}

type IfExistsReq struct {
	PhoneNumber string
}
type IfAdminExistsReq struct {
	PhoneNumber string
	Email       string
}

type IfExistsResp struct {
	IsExistsReq bool
}

type ChangeUserPasswordReq struct {
	PhoneNumber string
	Password    string
}

type ChangeAdminPasswordReq struct {
	Email       string
	PhoneNumber string
	Password    string
}

type ChangePasswordResp struct {
	Status bool
}

type ChangeAdminPasswordResp struct {
	Status bool
}

type UpdateRefreshTokenReq struct {
	Id           string
	RefreshToken string
}

type UpdateRefreshTokenResp struct {
	Status bool
}
