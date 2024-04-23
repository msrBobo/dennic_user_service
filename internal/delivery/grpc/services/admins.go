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

type adminRPC struct {
	logger         *zap.Logger
	admin          usecase.AdminStorageI
	brokerProducer event.BrokerProducer
}

func NewAdminRPC(logger *zap.Logger, admin usecase.AdminStorageI,
	brokerProducer event.BrokerProducer) pb.AdminServiceServer {
	return &adminRPC{
		logger:         logger,
		admin:          admin,
		brokerProducer: brokerProducer,
	}
}

func (a adminRPC) Create(ctx context.Context, admin *pb.Admin) (*pb.Admin, error) {

	req := entity.Admin{
		Id:            admin.Id,
		AdminOrder:    admin.AdminOrder,
		Role:          admin.Role,
		FirstName:     admin.FirstName,
		LastName:      admin.LastName,
		BirthDate:     admin.BirthDate,
		PhoneNumber:   admin.PhoneNumber,
		Email:         admin.Email,
		Password:      admin.Password,
		Gender:        admin.Gender,
		Salary:        admin.Salary,
		Biography:     admin.Biography,
		StartWorkYear: admin.StartWorkYear,
		EndWorkYear:   admin.EndWorkYear,
		WorkYears:     admin.WorkYears,
		RefreshToken:  admin.RefreshToken,
		CreatedAt:     time.Now(),
	}
	AdminId, err := a.admin.Create(ctx, &req)
	if err != nil {
		a.logger.Error("Create admin error", zap.Error(err))
		return nil, err
	}

	Params := make(map[string]string)
	Params["id"] = AdminId
	resp, err := a.admin.Get(ctx, Params)
	if err != nil {
		a.logger.Error("Create admin error", zap.Error(err))
		return nil, err
	}

	return &pb.Admin{
		Id:            resp.Id,
		AdminOrder:    resp.AdminOrder,
		Role:          resp.Role,
		FirstName:     resp.FirstName,
		LastName:      resp.LastName,
		BirthDate:     resp.BirthDate,
		PhoneNumber:   resp.PhoneNumber,
		Email:         resp.Email,
		Password:      resp.Password,
		Gender:        resp.Gender,
		Salary:        resp.Salary,
		Biography:     resp.Biography,
		StartWorkYear: resp.StartWorkYear,
		EndWorkYear:   resp.EndWorkYear,
		WorkYears:     resp.WorkYears,
		RefreshToken:  resp.RefreshToken,
		CreatedAt:     resp.CreatedAt.String(),
	}, nil
}

func (a adminRPC) Get(ctx context.Context, id *pb.GetAdminReqById) (*pb.Admin, error) {

	reqMap := make(map[string]string)
	reqMap["id"] = id.AdminId

	resp, err := a.admin.Get(ctx, reqMap)

	if err != nil {
		a.logger.Error("get admin error", zap.Error(err))
		return nil, err
	}

	return &pb.Admin{
		Id:            resp.Id,
		AdminOrder:    resp.AdminOrder,
		Role:          resp.Role,
		FirstName:     resp.FirstName,
		LastName:      resp.LastName,
		BirthDate:     resp.BirthDate,
		PhoneNumber:   resp.PhoneNumber,
		Email:         resp.Email,
		Password:      resp.Password,
		Gender:        resp.Gender,
		Salary:        resp.Salary,
		Biography:     resp.Biography,
		StartWorkYear: resp.StartWorkYear,
		EndWorkYear:   resp.EndWorkYear,
		WorkYears:     resp.WorkYears,
		RefreshToken:  resp.RefreshToken,
		CreatedAt:     resp.CreatedAt.String(),
		UpdatedAt:     resp.UpdatedAt.String(),
	}, nil
}

func (a adminRPC) ListAdmins(ctx context.Context, req *pb.ListAdminsReq) (*pb.ListAdminsResp, error) {

	resp, err := a.admin.List(ctx, req.Limit, req.Offset, req.Filter)

	if err != nil {
		a.logger.Error("get all admin error", zap.Error(err))
		return nil, err
	}

	var admins pb.ListAdminsResp

	for _, in := range resp {
		admins.Admins = append(admins.Admins, &pb.Admin{
			Id:            in.Id,
			AdminOrder:    in.AdminOrder,
			Role:          in.Role,
			FirstName:     in.FirstName,
			LastName:      in.LastName,
			BirthDate:     in.BirthDate,
			PhoneNumber:   in.PhoneNumber,
			Email:         in.Email,
			Password:      in.Password,
			Gender:        in.Gender,
			Salary:        in.Salary,
			Biography:     in.Biography,
			StartWorkYear: in.StartWorkYear,
			EndWorkYear:   in.EndWorkYear,
			WorkYears:     in.WorkYears,
			RefreshToken:  in.RefreshToken,
			CreatedAt:     in.CreatedAt.String(),
			UpdatedAt:     in.UpdatedAt.String(),
		})
	}

	return &admins, nil
}

func (a adminRPC) Update(ctx context.Context, admin *pb.Admin) (*pb.Admin, error) {

	req := entity.Admin{
		Id:            admin.Id,
		AdminOrder:    admin.AdminOrder,
		Role:          admin.Role,
		FirstName:     admin.FirstName,
		LastName:      admin.FirstName,
		BirthDate:     admin.Biography,
		PhoneNumber:   admin.PhoneNumber,
		Email:         admin.Email,
		Password:      admin.Password,
		Gender:        admin.Gender,
		Salary:        admin.Salary,
		Biography:     admin.Biography,
		StartWorkYear: admin.StartWorkYear,
		EndWorkYear:   admin.EndWorkYear,
		WorkYears:     admin.WorkYears,
		RefreshToken:  admin.RefreshToken,
		UpdatedAt:     time.Now(),
	}

	err := a.admin.Update(ctx, &req)

	if err != nil {
		a.logger.Error("update admin error", zap.Error(err))
		return nil, err
	}

	return &pb.Admin{
		Id:            req.Id,
		AdminOrder:    req.AdminOrder,
		Role:          req.Role,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		BirthDate:     req.BirthDate,
		PhoneNumber:   req.PhoneNumber,
		Email:         req.Email,
		Password:      req.Password,
		Gender:        req.Gender,
		Salary:        req.Salary,
		Biography:     req.Biography,
		StartWorkYear: req.StartWorkYear,
		EndWorkYear:   req.EndWorkYear,
		WorkYears:     req.WorkYears,
		RefreshToken:  req.RefreshToken,
		CreatedAt:     req.CreatedAt.String(),
		UpdatedAt:     req.UpdatedAt.String(),
	}, nil
}

func (a adminRPC) Delete(ctx context.Context, id *pb.DeleteAdminReq) (resp *emptypb.Empty, err error) {

	reqMap := make(map[string]string)
	reqMap["id"] = id.AdminId

	err = a.admin.Delete(ctx, id.AdminId)
	if err != nil {
		a.logger.Error("delete admin error", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

func (a adminRPC) CheckField(ctx context.Context, req *pb.CheckAdminFieldReq) (*pb.CheckAdminFieldResp, error) {

	reqAdmin := entity.CheckFieldReq{
		Value: req.Value,
		Field: req.Field,
	}

	resp, err := a.admin.CheckField(ctx, &reqAdmin)
	if err != nil {
		a.logger.Error("delete admin error", zap.Error(err))
		return nil, err
	}
	response := &pb.CheckAdminFieldResp{
		Status: resp.Status,
	}

	return response, nil
}

func (a adminRPC) IfExists(ctx context.Context, phone *pb.IfAdminExistsReq) (resp *pb.IfAdminExistsResp, err error) {

	req := entity.IfAdminExistsReq{
		PhoneNumber: phone.PhoneNumber,
		Email:       phone.Email,
	}

	entityResp, err := a.admin.IfExists(ctx, &req)

	if err != nil {
		a.logger.Error("delete admin error", zap.Error(err))
		return nil, err
	}

	resp = &pb.IfAdminExistsResp{
		IsExists: entityResp.IsExistsReq,
	}

	return resp, nil
}

func (a adminRPC) ChangePassword(ctx context.Context, phone *pb.ChangeAdminPasswordReq) (resp *pb.ChangeAdminPasswordResp, err error) {

	req := entity.ChangeAdminPasswordReq{
		Email:       phone.Email,
		PhoneNumber: phone.PhoneNumber,
		Password:    phone.Password,
	}
	status, err := a.admin.ChangePassword(ctx, &req)
	if err != nil {
		a.logger.Error("delete admin error", zap.Error(err))
		return nil, err
	}
	resp = &pb.ChangeAdminPasswordResp{
		Status: status.Status,
	}

	return resp, nil
}

func (a adminRPC) UpdateRefreshToken(ctx context.Context, id *pb.UpdateRefreshTokenAdminReq) (resp *pb.UpdateRefreshTokenAdminResp, err error) {
	req := entity.UpdateRefreshTokenReq{
		Id:           id.Id,
		RefreshToken: id.RefreshToken,
	}
	status, err := a.admin.UpdateRefreshToken(ctx, &req)
	if err != nil {
		a.logger.Error("delete admin error", zap.Error(err))
		return nil, err
	}

	resp = &pb.UpdateRefreshTokenAdminResp{
		Status: status.Status,
	}

	return resp, nil
}
