package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	fmt.Println("----")
	fmt.Println("subject:", subject, "description:", description)
	fmt.Println("----")

	// insert文にプレースホルダーで引数埋め込み
	result, err := s.db.ExecContext(ctx, insert, subject, description)

	fmt.Println("-------")
	fmt.Println(result)
	fmt.Println("-------")

	if err != nil {
		fmt.Println("--------")
		fmt.Println("insertでエラー発生")
		fmt.Println(err)
		fmt.Println("--------")

		return nil, err
	}

	// 自動挿入されたIDを取得する
	lastInsertId, err := result.LastInsertId()

	if err != nil {
		fmt.Println("自動挿入されたIDが取得できない")
		return nil, err
	}

	// select文にプレースホルダーで引数埋め込み
	var Todo model.TODO
	row := s.db.QueryRowContext(ctx, confirm, lastInsertId)

	// レコードを構造体に適合させる
	// TODO:&なかったらどうなるのか見てみたい
	row.Scan(&Todo.ID, &Todo.Subject, &Todo.Description, &Todo.CreatedAt, &Todo.UpdatedAt)

	fmt.Println("-----------------")
	fmt.Println(Todo)
	fmt.Println(&Todo)
	fmt.Println("-----------------")
	return &Todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	return nil, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	return nil, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
