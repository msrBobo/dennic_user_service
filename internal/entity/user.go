package entity

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
	CreatedAt    string
	UpdatedAt    string
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

type IfExistsResp struct {
	IsExistsReq bool
}

type ChangeUserPasswordReq struct {
	PhoneNumber string
	Password    string
}

type ChangeUserPasswordResp struct {
	Status bool
}

type UpdateRefreshTokenReq struct {
	UserId       string
	RefreshToken string
}

type UpdateRefreshTokenResp struct {
	Status bool
}
