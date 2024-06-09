package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

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

	// insert文にプレースホルダーで引数埋め込み
	result, err := s.db.ExecContext(ctx, insert, subject, description)

	if err != nil {
		fmt.Println("insertでエラー発生")
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
	rowError := row.Scan(&Todo.Subject, &Todo.Description, &Todo.CreatedAt, &Todo.UpdatedAt)

	// TODOにlastInsertIdを代入する
	Todo.ID = int64(lastInsertId)

	if rowError != nil {
		fmt.Println("select文でエラーが発生した。")
		return nil, rowError
	}

	return &Todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	var Todos []*model.TODO

	// prevIDあるかないか判定はとりあえず0でやる.引数渡さないの無理じゃない？
	if prevID == 0 {

		rows, err := s.db.QueryContext(ctx, read, size)

		if err != nil {
			fmt.Println("prevId:0のselect文でエラーが発生した。")
			return nil, err
		}

		for rows.Next() {
			var id int64
			var (
				subject, description string
			)
			var (
				createdAt, updatedAt time.Time
			)

			err := rows.Scan(&id, &subject, &description, &createdAt, &updatedAt)

			if err != nil {
				fmt.Println("prevId:0のscanでエラーが発生した")
				return nil, err
			}

			Todo := model.TODO{
				ID:          id,
				Subject:     subject,
				Description: description,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			}

			Todos = append(Todos, &Todo)
		}

		fmt.Println(Todos)

		// 何もない場合は空スライスを代入する(nilスライスにしない)
		if len(Todos) == 0 {
			Todos = make([]*model.TODO, 0)
		}

		return Todos, nil
	}

	// 複数レコード取得時
	rows, err := s.db.QueryContext(ctx, readWithID, prevID, size)

	if err != nil {
		fmt.Println("selectでエラー発生")
		return nil, err
	}

	for rows.Next() {
		var id int64
		var (
			subject, description string
		)
		var (
			createdAt, updatedAt time.Time
		)

		err := rows.Scan(&id, &subject, &description, &createdAt, &updatedAt)

		if err != nil {
			fmt.Println("select文でエラーが発生した")
			return nil, err
		}

		Todo := model.TODO{
			ID:          id,
			Subject:     subject,
			Description: description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		Todos = append(Todos, &Todo)
	}

	// 何もない場合は空スライスを代入する(nilスライスにしない)
	if len(Todos) == 0 {
		Todos = make([]*model.TODO, 0)
	}

	return Todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// 先にupdateクエリを流す
	result, err := s.db.ExecContext(ctx, update, subject, description, id)

	if err != nil {
		fmt.Println("insertでエラー発生")
		return nil, err
	}

	// result見てrowの数取得したい.0だったらエラー返す
	rows, _ := result.RowsAffected()
	if rows == 0 {
		// Errorってメソッド名ならエラー発生時に呼ばれる？
		var errNotFound model.ErrNotFound
		return nil, &errNotFound
	}

	// updateされた後のレコードを取得してレスポンスとして返す
	var Todo model.TODO
	row := s.db.QueryRowContext(ctx, confirm, id)

	rowError := row.Scan(&Todo.Subject, &Todo.Description, &Todo.CreatedAt, &Todo.UpdatedAt)
	if rowError != nil {
		fmt.Println("select文でエラーが発生した。")
		return nil, rowError
	}

	Todo.ID = id

	return &Todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	// 最初の1個分減らしとく
	placeholder := strings.Repeat(", ?", len(ids)-1)
	deleteQuery := fmt.Sprintf(deleteFmt, placeholder)

	// delete文にプレースホルダーを埋め込む
	var idsInterface []interface{}
	for _, id := range ids {
		idsInterface = append(idsInterface, id)
	}

	// 配列を...で展開する
	rows, err := s.db.ExecContext(ctx, deleteQuery, idsInterface...)

	if err != nil {
		fmt.Println("delete文でエラーが発生した")
		return err
	}

	rowsAffected, _ := rows.RowsAffected()

	if rowsAffected == 0 {
		fmt.Println("レコードが削除されなかった場合発火する")
		var errNotFound model.ErrNotFound
		return &errNotFound
	}

	return nil
}
