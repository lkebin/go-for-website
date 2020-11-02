package post

import (
	"buhaoyong/pkg/db"
	"buhaoyong/pkg/db/helper"
	"buhaoyong/pkg/service/post/model"
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, m *model.Post) (*model.Post, error)
	Update(ctx context.Context, m *model.Post) error
	FindByID(ctx context.Context, id int32) (*model.Post, error)
	Delete(ctx context.Context, id int32) error
}

type repositoryImpl struct {
	db *db.DB
}

func New(dbc *db.DB) Repository {
	return &repositoryImpl{db: dbc}
}

func (r *repositoryImpl) Create(ctx context.Context, m *model.Post) (*model.Post, error) {
	columns, placeholders, values, err := helper.BuildForInsert(m)
	if err != nil {
		return nil, err
	}

	result, err := r.db.ExecContext(ctx, fmt.Sprintf(`insert into post (%s) values(%s)`, columns, placeholders), values...)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, int32(id))
}

func (r *repositoryImpl) Update(ctx context.Context, m *model.Post) error {
	columns, values, err := helper.BuildForUpdate(m)
	if err != nil {
		return err
	}

	values = append(values, m.ID)
	_, err = r.db.ExecContext(ctx, fmt.Sprintf(`update post set %s where id=?`, columns), values...)
	return err
}

func (r *repositoryImpl) FindByID(ctx context.Context, id int32) (*model.Post, error) {
	var m model.Post
	if err := r.db.QueryRowxContext(ctx, `select id,title,content from post where id=?`, id).StructScan(&m); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &m, nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id int32) error {
	_, err := r.db.ExecContext(ctx, `delete from post where id=?`, id)
	return err
}
