package repository

import (
	"context"
	"dennic_user_service/internal/entity"
)

type AdminStorageI interface {
	Create(ctx context.Context, admin *entity.Admin) error
	Get(ctx context.Context, params map[string]string) (*entity.Admin, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Admin, error)
	Update(ctx context.Context, kyc *entity.Admin) error
	Delete(ctx context.Context, id string) error
	CheckField(ctx context.Context, req *entity.CheckFieldReq) (*entity.CheckFieldResp, error)
	IfExists(ctx context.Context, req *entity.IfExistsReq) (*entity.IfExistsResp, error)
	ChangePassword(ctx context.Context, req *entity.ChangeUserPasswordReq) (*entity.ChangePasswordResp, error)
	UpdateRefreshToken(ctx context.Context, req *entity.UpdateRefreshTokenReq) (*entity.UpdateRefreshTokenResp, error)
}
