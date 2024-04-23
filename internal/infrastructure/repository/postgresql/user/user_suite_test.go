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

type UserReposisitoryTestSuite struct {
	suite.Suite
	Config     *config.Config
	DB         *postgres.PostgresDB
	repo       repo.UserStorageI
	ctxTimeout time.Duration
}

func NewUserService(ctxTimeout time.Duration, repo repo.UserStorageI, config *config.Config) UserReposisitoryTestSuite {
	return UserReposisitoryTestSuite{
		Config:     config,
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

// test func
func (s *UserReposisitoryTestSuite) TestUserCRUD() {

	config := config.New()

	db, err := postgres.New(config)
	if err != nil {
		s.T().Fatal("Error initializing database connection:", err)
	}

	s.DB = db

	userRepo := NewUserRepo(s.DB)
	ctx := context.Background()

	// struct for create user
	user := entity.User{
		UserOrder:    101,
		FirstName:    "firstname",
		LastName:     "lastname",
		BirthDate:    "2000-08-30",
		PhoneNumber:  "+998994767316",
		Password:     "testpassword",
		Gender:       "male",
		RefreshToken: "testrefreshtoken",
		CreatedAt:    time.Now().UTC(),
	}
	// uuid generating
	user.Id = uuid.New().String()

	updUser := entity.User{
		Id:          user.Id,
		FirstName:   "updfirstname",
		LastName:    "updlastname",
		BirthDate:   "2000-07-20",
		PhoneNumber: "+998934767316",
		Password:    "updtestpassword",
		Gender:      "male",
		UpdatedAt:   time.Now(),
	}

	// check create user method
	err = userRepo.Create(ctx, &user)
	s.Suite.NoError(err)
	Params := make(map[string]string)
	Params["id"] = user.Id

	// check get user method
	getUser, err := userRepo.Get(ctx, Params)
	s.Suite.NoError(err)
	s.Suite.NotNil(getUser)
	s.Suite.Equal(getUser.Id, user.Id)
	s.Suite.Equal(getUser.FirstName, user.FirstName)
	s.Suite.Equal(getUser.PhoneNumber, user.PhoneNumber)

	// check update user method
	err = userRepo.Update(ctx, &updUser)
	s.Suite.NoError(err)
	updGetUser, err := userRepo.Get(ctx, Params)
	s.Suite.NoError(err)
	s.Suite.NotNil(updGetUser)
	s.Suite.Equal(updGetUser.Id, updUser.Id)
	s.Suite.Equal(updGetUser.FirstName, updUser.FirstName)
	s.Suite.Equal(updGetUser.PhoneNumber, updUser.PhoneNumber)

	// check getAllUsers method
	getAllUsers, err := userRepo.List(ctx, 5, 1, nil)
	s.Suite.NoError(err)
	s.Suite.NotNil(getAllUsers)
	req := entity.CheckFieldReq{
		Value: updUser.PhoneNumber,
		Field: "phone_number",
	}

	// check CheckField user method
	result, err := userRepo.CheckField(ctx, &req)
	s.Suite.NoError(err)
	s.Suite.NotNil(updGetUser)
	s.Suite.Equal(result.Status, true)

	// check IfExists user method
	if_exists_req := entity.IfExistsReq{
		PhoneNumber: updUser.PhoneNumber,
	}
	status, err := userRepo.IfExists(ctx, &if_exists_req)
	s.Suite.NoError(err)
	s.Suite.NotNil(updGetUser)
	s.Suite.Equal(status.IsExistsReq, true)

	// check ChangePassword user method
	change_password_req := entity.ChangeUserPasswordReq{
		PhoneNumber: updUser.PhoneNumber,
		Password:    "new_password",
	}
	resp_change_password, err := userRepo.ChangePassword(ctx, &change_password_req)
	s.Suite.NoError(err)
	s.Suite.NotNil(resp_change_password)
	s.Suite.Equal(resp_change_password.Status, true)

	// check UpdateRefreshToken user method
	req_update_refresh_token := entity.UpdateRefreshTokenReq{
		Id:           updUser.Id,
		RefreshToken: "new_refresh_token",
	}
	resp_update_refresh_token, err := userRepo.UpdateRefreshToken(ctx, &req_update_refresh_token)
	s.Suite.NoError(err)
	s.Suite.NotNil(resp_update_refresh_token)
	s.Suite.Equal(resp_update_refresh_token.Status, true)

	//check delete user method
	err = userRepo.Delete(ctx, user.Id)
	s.Suite.NoError(err)

}

func TestExampleUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserReposisitoryTestSuite))
}
