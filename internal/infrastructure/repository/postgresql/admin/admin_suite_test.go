package postgresql

import (
	"context"
	"dennic_user_service/internal/entity"
	repo "dennic_user_service/internal/infrastructure/repository"
	"dennic_user_service/internal/pkg/config"
	"dennic_user_service/internal/pkg/postgres"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/suite"
)

type AdminReposisitoryTestSuite struct {
	suite.Suite
	Config     *config.Config
	DB         *postgres.PostgresDB
	repo       repo.AdminStorageI
	ctxTimeout time.Duration
}

func NewAdminService(ctxTimeout time.Duration, repo repo.AdminStorageI, config *config.Config) AdminReposisitoryTestSuite {
	return AdminReposisitoryTestSuite{
		Config:     config,
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

// test func
func (s *AdminReposisitoryTestSuite) TestAdminCRUD() {

	config := config.New()

	db, err := postgres.New(config)
	if err != nil {
		s.T().Fatal("Error initializing database connection:", err)
	}

	s.DB = db

	adminRepo := NewAdminRepo(s.DB)
	ctx := context.Background()

	// struct for create admin
	admin := entity.Admin{
		AdminOrder:    777,
		Role:          "admin",
		FirstName:     "testdata",
		LastName:      "testdata",
		BirthDate:     "2000-08-30",
		PhoneNumber:   "testdata",
		Email:         "testdata",
		Password:      "testdata",
		Gender:        "male",
		Salary:        777.7,
		Biography:     "testdata",
		StartWorkYear: "2000-08-30",
		EndWorkYear:   "2000-08-30",
		WorkYears:     777,
		RefreshToken:  "testdata",
		CreatedAt:     time.Now().UTC(),
	}
	// uuid generating
	admin.Id = uuid.New().String()

	updAdmin := entity.Admin{
		Id:            admin.Id,
		AdminOrder:    888,
		Role:          "superadmin",
		FirstName:     "updtestdata",
		LastName:      "updtestdata",
		BirthDate:     "2000-07-20",
		PhoneNumber:   "updtestdata",
		Email:         "updtestdata",
		Password:      "updtestdata",
		Gender:        "male",
		Salary:        888,
		Biography:     "updtestdata",
		StartWorkYear: "2000-08-30",
		EndWorkYear:   "2000-08-30",
		WorkYears:     888,
		RefreshToken:  "updtestdata",
		UpdatedAt:     time.Now(),
	}
	_ = updAdmin.UpdatedAt
	// check create admin method
	err = adminRepo.Create(ctx, &admin)
	s.Suite.NoError(err)
	Params := make(map[string]string)
	Params["id"] = admin.Id

	// check get admin method
	getAdmin, err := adminRepo.Get(ctx, Params)
	s.Suite.NoError(err)
	s.Suite.NotNil(getAdmin)
	s.Suite.Equal(getAdmin.Id, admin.Id)
	s.Suite.Equal(getAdmin.FirstName, admin.FirstName)
	s.Suite.Equal(getAdmin.PhoneNumber, admin.PhoneNumber)

	// check update admin method
	err = adminRepo.Update(ctx, &updAdmin)
	s.Suite.NoError(err)
	updGetAdmin, err := adminRepo.Get(ctx, Params)
	s.Suite.NoError(err)
	s.Suite.NotNil(updGetAdmin)
	s.Suite.Equal(updGetAdmin.Id, updAdmin.Id)
	s.Suite.Equal(updGetAdmin.FirstName, updAdmin.FirstName)
	s.Suite.Equal(updGetAdmin.PhoneNumber, updAdmin.PhoneNumber)
	s.Suite.Equal(updGetAdmin.Role, updAdmin.Role)

	// // check getAllAdmins method
	getAllAdmins, err := adminRepo.List(ctx, 5, 1, nil)
	s.Suite.NoError(err)
	s.Suite.NotNil(getAllAdmins)
	req := entity.CheckFieldReq{
		Value: updAdmin.Email,
		Field: "email",
	}

	// // check CheckField admin method
	result, err := adminRepo.CheckField(ctx, &req)
	s.Suite.NoError(err)
	s.Suite.NotNil(updGetAdmin)
	s.Suite.Equal(result.Status, true)

	// // check IfExists user method
	if_exists_req := entity.IfAdminExistsReq{
		PhoneNumber: updAdmin.PhoneNumber,
		Email:       updAdmin.Email,
	}
	status, err := adminRepo.IfExists(ctx, &if_exists_req)
	s.Suite.NoError(err)
	s.Suite.NotNil(updGetAdmin)
	s.Suite.Equal(status.IsExistsReq, true)

	// check ChangePassword user method for PhoneNumber change
	change_password_req := entity.ChangeAdminPasswordReq{
		PhoneNumber: updAdmin.PhoneNumber,
		Password:    "new_password",
	}

	resp_change_password, err := adminRepo.ChangePassword(ctx, &change_password_req)
	s.Suite.NoError(err)
	s.Suite.NotNil(resp_change_password)
	s.Suite.Equal(resp_change_password.Status, true)

    // check ChangePassword user method for Email change
	change_password_req_2 := entity.ChangeAdminPasswordReq{
		Email: updAdmin.Email,
		Password:    "new_password",
	}
	resp_change_password_2, err := adminRepo.ChangePassword(ctx, &change_password_req_2)
	s.Suite.NoError(err)
	s.Suite.NotNil(resp_change_password_2)
	s.Suite.Equal(resp_change_password_2.Status, true)

	// // check UpdateRefreshToken admin method
	req_update_refresh_token := entity.UpdateRefreshTokenReq{
		Id:       updAdmin.Id,
		RefreshToken: "new_refresh_token",
	}
	resp_update_refresh_token, err := adminRepo.UpdateRefreshToken(ctx, &req_update_refresh_token)
	s.Suite.NoError(err)
	s.Suite.NotNil(resp_update_refresh_token)
	s.Suite.Equal(resp_update_refresh_token.Status, true)

	// // check delete admin method
	err = adminRepo.Delete(ctx, admin.Id)
	s.Suite.NoError(err)

}

func TestExampleAdminTestSuite(t *testing.T) {
	suite.Run(t, new(AdminReposisitoryTestSuite))
}
