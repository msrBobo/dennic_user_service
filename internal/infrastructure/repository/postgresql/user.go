package postgresql

import (
	"context"
	"dennic_user_service/internal/entity"
	"dennic_user_service/internal/pkg/otlp"
	"dennic_user_service/internal/pkg/postgres"
	"fmt"

	"github.com/Masterminds/squirrel"
)

const (
	userTableName      = "users"
	userServiceName    = "userService"
	userSpanRepoPrefix = "userRepo"
)

type userRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewUserRepo(db *postgres.PostgresDB) *userRepo {
	return &userRepo{
		tableName: userTableName,
		db:        db,
	}
}

func (p *userRepo) userSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.
		Select(
			"id",
			"user_order",
			"first_name",
			"last_name",
			"birth_date",
			"phone_number",
			"password",
			"gender",
			"created_at",
		).From(p.tableName)
}

func (p userRepo) Create(ctx context.Context, user *entity.User) error {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Create")
	defer span.End()
	data := map[string]any{
		"id":            user.Id,
		"user_order":    user.UserOrder,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"birth_date":    user.BirthDate,
		"phone_number":  user.PhoneNumber,
		"password":      user.Password,
		"gender":        user.Gender,
		"refresh_token": user.RefreshToken,
		"created_at":    user.CreatedAt,
	}

	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

func (p userRepo) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Get")
	defer span.End()

	var (
		user entity.User
	)

	queryBuilder := p.userSelectQueryPrefix()

	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		}
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}
	if err = p.db.QueryRow(ctx, query, args...).Scan(
		&user.Id,
		&user.UserOrder,
		&user.FirstName,
		&user.LastName,
		&user.BirthDate,
		&user.PhoneNumber,
		&user.Password,
		&user.Gender,
		&user.CreatedAt,
	); err != nil {
		return nil, p.db.Error(err)
	}

	return &user, nil
}

func (p userRepo) List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.User, error) {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"List")
	defer span.End()

	var (
		users []*entity.User
	)
	queryBuilder := p.userSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(limit).Offset(offset)
	}

	for key, value := range filter {
		if key == "id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
			continue
		}
		if key == "created_at" {
			queryBuilder = queryBuilder.Where("created_at=?", value)
			continue
		}
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "list"))
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()
	users = make([]*entity.User, 0)
	for rows.Next() {
		var user entity.User
		if err = rows.Scan(
			&user.Id,
			&user.UserOrder,
			&user.FirstName,
			&user.LastName,
			&user.BirthDate,
			&user.PhoneNumber,
			&user.Password,
			&user.Gender,
			&user.CreatedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}

		users = append(users, &user)
	}

	return users, nil
}

func (p userRepo) Update(ctx context.Context, user *entity.User) error {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Update")
	defer span.End()
	clauses := map[string]any{
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"birth_date":   user.BirthDate,
		"phone_number": user.PhoneNumber,
		"password":     user.Password,
		"gender":       user.Gender,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", user.Id)).
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (p userRepo) Delete(ctx context.Context, id string) error {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Delete")
	defer span.End()

	sqlStr, args, err := p.db.Sq.Builder.
		Delete(p.tableName).
		Where(p.db.Sq.Equal("id", id)).
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" delete")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (p *userRepo) CheckField(ctx context.Context, req *entity.CheckFieldReq) (*entity.CheckFieldResp, error) {
	query := fmt.Sprintf(`SELECT count(1) FROM users WHERE %s = $1 AND deleted_at IS NULL`, req.Field)

	var isExists int

	row := p.db.QueryRow(ctx, query, req.Value)
	if err := row.Scan(&isExists); err != nil {
		return nil, err
	}

	if isExists == 1 {
		return &entity.CheckFieldResp{
			Status: true,
		}, nil
	}

	return &entity.CheckFieldResp{
		Status: false,
	}, nil
}

func (p *userRepo) IfExists(ctx context.Context, req *entity.IfExistsReq) (resp *entity.IfExistsResp, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1 AND deleted_at IS NULL)`

	var exists bool
	row := p.db.QueryRow(ctx, query, req.PhoneNumber)
	if err = row.Scan(&exists); err != nil {
		return nil, err
	}
	resp = &entity.IfExistsResp{
		IsExistsReq: exists,
	}

	return resp, nil
}

func (p *userRepo) ChangePassword(ctx context.Context, req *entity.ChangeUserPasswordReq) (*entity.ChangeUserPasswordResp, error) {
	query := `UPDATE users SET password = $1 WHERE phone_number = $2 AND deleted_at IS NULL`
	resp, err := p.db.Exec(ctx, query, req.Password, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if resp.RowsAffected() == 0 {
		return &entity.ChangeUserPasswordResp{Status: false}, nil
	}
	return &entity.ChangeUserPasswordResp{Status: true}, nil
}

func (p *userRepo) UpdateRefreshToken(ctx context.Context, req *entity.UpdateRefreshTokenReq) (*entity.UpdateRefreshTokenResp, error) {
	query := `UPDATE users SET refresh_token = $1 WHERE id = $2 AND deleted_at IS NULL`

	resp, err := p.db.Exec(ctx, query, req.RefreshToken, req.UserId)
	if err != nil {
		return nil, err
	}
	if resp.RowsAffected() == 0 {
		return &entity.UpdateRefreshTokenResp{Status: false}, nil
	}

	return &entity.UpdateRefreshTokenResp{Status: true}, nil
}
