package services

import (
	"context"
	pb "dennic_user_service/genproto/user_service"
	"dennic_user_service/internal/entity"
	"dennic_user_service/internal/usecase"
	"dennic_user_service/internal/usecase/event"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userRPC struct {
	logger         *zap.Logger
	user           usecase.UserStorageI
	brokerProducer event.BrokerProducer
}

func NewUserRPC(logger *zap.Logger, user usecase.UserStorageI,
	brokerProducer event.BrokerProducer) pb.UserServiceServer {
	return &userRPC{
		logger:         logger,
		user:           user,
		brokerProducer: brokerProducer,
	}
}

func (u userRPC) Create(ctx context.Context, user *pb.User) (*pb.User, error) {

	req := entity.User{
		Id:           user.Id,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		BirthDate:    user.BirthDate,
		PhoneNumber:  user.PhoneNumber,
		Password:     user.Password,
		Gender:       user.Gender,
		RefreshToken: user.RefreshToken,
		CreatedAt:    time.Now(),
	}
	UserId, err := u.user.Create(ctx, &req)
	if err != nil {
		u.logger.Error("Create user error", zap.Error(err))
		return nil, err
	}

	Params := make(map[string]string)
	Params["id"] = UserId
	resp, err := u.user.Get(ctx, Params)
	if err != nil {
		u.logger.Error("Create user error", zap.Error(err))
		return nil, err
	}

	return &pb.User{
		Id:           resp.Id,
		UserOrder:    resp.UserOrder,
		FirstName:    resp.FirstName,
		LastName:     resp.LastName,
		BirthDate:    resp.BirthDate,
		PhoneNumber:  resp.PhoneNumber,
		Password:     resp.Password,
		Gender:       resp.Gender,
		RefreshToken: resp.RefreshToken,
		CreatedAt:    resp.CreatedAt.String(),
	}, nil
}

func (u userRPC) Get(ctx context.Context, id *pb.GetUserReqById) (*pb.User, error) {

	reqMap := make(map[string]string)
	reqMap["id"] = id.UserId

	resp, err := u.user.Get(ctx, reqMap)

	if err != nil {
		u.logger.Error("get user error", zap.Error(err))
		return nil, err
	}

	return &pb.User{
		Id:           resp.Id,
		UserOrder:    resp.UserOrder,
		FirstName:    resp.FirstName,
		LastName:     resp.LastName,
		BirthDate:    resp.BirthDate,
		PhoneNumber:  resp.PhoneNumber,
		Password:     resp.Password,
		Gender:       resp.Gender,
		RefreshToken: resp.RefreshToken,
		CreatedAt:    resp.CreatedAt.String(),
		UpdatedAt:    resp.UpdatedAt.String(),
	}, nil
}

func (u userRPC) ListUsers(ctx context.Context, req *pb.ListUsersReq) (*pb.ListUsersResp, error) {

	resp, err := u.user.List(ctx, req.Limit, req.Offset, req.Filter)

	if err != nil {
		u.logger.Error("get all user error", zap.Error(err))
		return nil, err
	}

	var users pb.ListUsersResp

	for _, in := range resp {
		users.Users = append(users.Users, &pb.User{
			Id:           in.Id,
			UserOrder:    in.UserOrder,
			FirstName:    in.FirstName,
			LastName:     in.LastName,
			BirthDate:    in.BirthDate,
			PhoneNumber:  in.PhoneNumber,
			Password:     in.Password,
			Gender:       in.Gender,
			RefreshToken: in.RefreshToken,
			CreatedAt:    in.CreatedAt.String(),
			UpdatedAt:    in.UpdatedAt.String(),
		})
	}

	return &users, nil
}

func (u userRPC) Update(ctx context.Context, user *pb.User) (*pb.User, error) {

	req := entity.User{
		Id:           user.Id,
		UserOrder:    user.UserOrder,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		BirthDate:    user.BirthDate,
		PhoneNumber:  user.PhoneNumber,
		Password:     user.Password,
		Gender:       user.Gender,
		RefreshToken: user.RefreshToken,
	}

	err := u.user.Update(ctx, &req)

	if err != nil {
		u.logger.Error("update user error", zap.Error(err))
		return nil, err
	}

	return &pb.User{
		Id:           req.Id,
		UserOrder:    req.UserOrder,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BirthDate:    req.BirthDate,
		PhoneNumber:  req.PhoneNumber,
		Password:     req.Password,
		Gender:       req.Gender,
		RefreshToken: req.RefreshToken,
		CreatedAt:    req.CreatedAt.String(),
		UpdatedAt:    req.UpdatedAt.String(),
	}, nil
}

func (u userRPC) Delete(ctx context.Context, id *pb.DeleteUserReq) (resp *emptypb.Empty, err error) {

	reqMap := make(map[string]string)
	reqMap["id"] = id.UserId

	err = u.user.Delete(ctx, id.UserId)
	if err != nil {
		u.logger.Error("delete user error", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

func (u userRPC) CheckField(ctx context.Context, req *pb.CheckFieldUserReq) (*pb.CheckFieldUserResp, error) {

	reqUser := entity.CheckFieldReq{
		Value: req.Value,
		Field: req.Field,
	}

	resp, err := u.user.CheckField(ctx, &reqUser)
	if err != nil {
		u.logger.Error("delete user error", zap.Error(err))
		return nil, err
	}
	response := &pb.CheckFieldUserResp{
		Status: resp.Status,
	}

	return response, nil
}

func (u userRPC) IfExists(ctx context.Context, phone *pb.IfUserExistsReq) (resp *pb.IfUserExistsResp, err error) {

	req := entity.IfExistsReq{
		PhoneNumber: phone.PhoneNumber,
	}

	entityResp, err := u.user.IfExists(ctx, &req)

	if err != nil {
		u.logger.Error("delete user error", zap.Error(err))
		return nil, err
	}

	resp = &pb.IfUserExistsResp{
		IsExists: entityResp.IsExistsReq,
	}

	return resp, nil
}

func (u userRPC) ChangePassword(ctx context.Context, phone *pb.ChangeUserPasswordReq) (resp *pb.ChangeUserPasswordResp, err error) {

	req := entity.ChangeUserPasswordReq{
		PhoneNumber: phone.PhoneNumber,
		Password:    phone.Password,
	}
	status, err := u.user.ChangePassword(ctx, &req)
	if err != nil {
		u.logger.Error("delete user error", zap.Error(err))
		return nil, err
	}
	resp = &pb.ChangeUserPasswordResp{
		Status: status.Status,
	}

	return resp, nil
}

func (u userRPC) UpdateRefreshToken(ctx context.Context, id *pb.UpdateRefreshTokenUserReq) (resp *pb.UpdateRefreshTokenUserResp, err error) {
	req := entity.UpdateRefreshTokenReq{
		Id:           id.Id,
		RefreshToken: id.RefreshToken,
	}
	status, err := u.user.UpdateRefreshToken(ctx, &req)
	if err != nil {
		u.logger.Error("delete user error", zap.Error(err))
		return nil, err
	}

	resp = &pb.UpdateRefreshTokenUserResp{
		Status: status.Status,
	}

	return resp, nil
}
