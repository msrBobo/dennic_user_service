package postgresql

import (
	"context"
	"database/sql"
	"dennic_user_service/internal/entity"
	"dennic_user_service/internal/pkg/otlp"
	"dennic_user_service/internal/pkg/postgres"
	"fmt"

	"github.com/Masterminds/squirrel"
)

const (
	adminTableName      = "admins"
	adminServiceName    = "adminService"
	adminSpanRepoPrefix = "adminRepo"
)

type adminRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewAdminRepo(db *postgres.PostgresDB) *adminRepo {
	return &adminRepo{
		tableName: adminTableName,
		db:        db,
	}
}

func (p *adminRepo) adminSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.
		Select(
			"id",
			"admin_order",
			"role",
			"first_name",
			"last_name",
			"birth_date",
			"phone_number",
			"email",
			"password",
			"gender",
			"salary",
			"biography",
			"start_work_year",
			"end_work_year",
			"work_years",
			"created_at",
			"updated_at",
		).From(p.tableName).
		Where("deleted_at IS NULL")
}

func (p adminRepo) Create(ctx context.Context, admin *entity.Admin) error {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"Create")
	defer span.End()
	data := map[string]any{
		"id":              admin.Id,
		"role":            admin.Role,
		"first_name":      admin.FirstName,
		"last_name":       admin.LastName,
		"birth_date":      admin.BirthDate,
		"phone_number":    admin.PhoneNumber,
		"email":           admin.Email,
		"password":        admin.Password,
		"gender":          admin.Gender,
		"salary":          admin.Salary,
		"biography":       admin.Biography,
		"start_work_year": admin.StartWorkYear,
		"end_work_year":   admin.EndWorkYear,
		"work_years":      admin.WorkYears,
		"refresh_token":   admin.RefreshToken,
		"created_at":      admin.CreatedAt,
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

func (p adminRepo) Get(ctx context.Context, params map[string]string) (*entity.Admin, error) {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"Get")
	defer span.End()

	var (
		admin entity.Admin
	)

	queryBuilder := p.adminSelectQueryPrefix()

	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		}
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}

	var (
		birthDate       sql.NullString
		updatedAt       sql.NullTime
		start_work_year sql.NullString
		end_work_year   sql.NullString
	)
	if err = p.db.QueryRow(ctx, query, args...).Scan(
		&admin.Id,
		&admin.AdminOrder,
		&admin.Role,
		&admin.FirstName,
		&admin.LastName,
		&birthDate,
		&admin.PhoneNumber,
		&admin.Email,
		&admin.Password,
		&admin.Gender,
		&admin.Salary,
		&admin.Biography,
		&start_work_year,
		&end_work_year,
		&admin.WorkYears,
		&admin.CreatedAt,
		&updatedAt,
	); err != nil {
		return nil, p.db.Error(err)
	}

	if updatedAt.Valid {
		admin.UpdatedAt = updatedAt.Time
	}
	if birthDate.Valid {
		admin.BirthDate = birthDate.String
	}
	if start_work_year.Valid {
		admin.StartWorkYear = start_work_year.String
	}
	if end_work_year.Valid {
		admin.EndWorkYear = end_work_year.String
	}

	return &admin, nil
}

func (p adminRepo) List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Admin, error) {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"List")
	defer span.End()

	var (
		admins []*entity.Admin
	)
	queryBuilder := p.adminSelectQueryPrefix()

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
	var (
		birthDate       sql.NullString
		updatedAt       sql.NullTime
		start_work_year sql.NullString
		end_work_year   sql.NullString
	)
	for rows.Next() {
		var admin entity.Admin
		if err = rows.Scan(
			&admin.Id,
			&admin.AdminOrder,
			&admin.Role,
			&admin.FirstName,
			&admin.LastName,
			&birthDate,
			&admin.PhoneNumber,
			&admin.Email,
			&admin.Password,
			&admin.Gender,
			&admin.Salary,
			&admin.Biography,
			&start_work_year,
			&end_work_year,
			&admin.WorkYears,
			&admin.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}

		if updatedAt.Valid {
			admin.UpdatedAt = updatedAt.Time
		}
		if birthDate.Valid {
			admin.BirthDate = birthDate.String
		}
		if start_work_year.Valid {
			admin.StartWorkYear = start_work_year.String
		}
		if end_work_year.Valid {
			admin.EndWorkYear = end_work_year.String
		}
		admins = append(admins, &admin)
	}

	return admins, nil
}

func (p *adminRepo) Update(ctx context.Context, admin *entity.Admin) error {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"Update")
	defer span.End()

	clauses := map[string]interface{}{
		"role":            admin.Role,
		"first_name":      admin.FirstName,
		"last_name":       admin.LastName,
		"birth_date":      admin.BirthDate,
		"phone_number":    admin.PhoneNumber,
		"email":           admin.Email,
		"password":        admin.Password,
		"gender":          admin.Gender,
		"salary":          admin.Salary,
		"biography":       admin.Biography,
		"start_work_year": admin.StartWorkYear,
		"end_work_year":   admin.EndWorkYear,
		"work_years":      admin.WorkYears,
		"updated_at":      admin.UpdatedAt,
	}

	updateBuilder := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", admin.Id))

	updateBuilder = updateBuilder.Where("deleted_at IS NULL")

	sqlStr, args, err := updateBuilder.ToSql()
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

func (p *adminRepo) Delete(ctx context.Context, id string) error {
	ctx, span := otlp.Start(ctx, adminServiceName, adminSpanRepoPrefix+"Delete")
	defer span.End()

	sqlStr := fmt.Sprintf(`
		UPDATE %s
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, p.tableName)

	commandTag, err := p.db.Exec(ctx, sqlStr, id)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (p *adminRepo) CheckField(ctx context.Context, req *entity.CheckFieldReq) (*entity.CheckFieldResp, error) {
	query := fmt.Sprintf(`
		SELECT count(1) 
			FROM admins WHERE %s = $1 AND 
			deleted_at IS NULL`, req.Field)

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

func (p *adminRepo) IfExists(ctx context.Context, req *entity.IfAdminExistsReq) (resp *entity.IfExistsResp, err error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM admins
			WHERE (phone_number = $1 AND email = $2)
			AND deleted_at IS NULL
		)
	`
	var exists bool
	row := p.db.QueryRow(ctx, query, req.PhoneNumber, req.Email)
	if err = row.Scan(&exists); err != nil {
		return nil, err
	}
	resp = &entity.IfExistsResp{
		IsExistsReq: exists,
	}

	return resp, nil
}

func (p *adminRepo) ChangePassword(ctx context.Context, req *entity.ChangeAdminPasswordReq) (*entity.ChangeAdminPasswordResp, error) {
	query := `
		UPDATE admins
		SET password = $1
		WHERE (email = $2 OR phone_number = $3)
		AND deleted_at IS NULL
	`

	resp, err := p.db.Exec(ctx, query, req.Password, req.Email, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	if resp.RowsAffected() == 0 {
		return &entity.ChangeAdminPasswordResp{Status: false}, nil
	}

	return &entity.ChangeAdminPasswordResp{Status: true}, nil
}

func (p *adminRepo) UpdateRefreshToken(ctx context.Context, req *entity.UpdateRefreshTokenReq) (*entity.UpdateRefreshTokenResp, error) {
	query := `
			UPDATE admins 
			SET refresh_token = $1 
			WHERE id = $2 AND 
			deleted_at IS NULL`

	resp, err := p.db.Exec(ctx, query, req.RefreshToken, req.Id)
	if err != nil {
		return nil, err
	}
	if resp.RowsAffected() == 0 {
		return &entity.UpdateRefreshTokenResp{Status: false}, nil
	}

	return &entity.UpdateRefreshTokenResp{Status: true}, nil
}
