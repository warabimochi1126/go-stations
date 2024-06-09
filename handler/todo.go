package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	_, _ = h.svc.CreateTODO(ctx, "", "")
	return &model.CreateTODOResponse{}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// r.ContentLengthでbodyの長さ取ってきてる
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		var createTodoReq model.CreateTODORequest
		json.Unmarshal(body, &createTodoReq)

		if createTodoReq.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			todo, _ := h.svc.CreateTODO(r.Context(), createTodoReq.Subject, createTodoReq.Description)

			var responseTodo model.CreateTODOResponse
			responseTodo.TODO = *todo

			err := json.NewEncoder(w).Encode(responseTodo)

			if err != nil {
				fmt.Println(err)
			}

		}
	} else if r.Method == "PUT" {
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		var updateTodoReq model.UpdateTODORequest
		json.Unmarshal(body, &updateTodoReq)

		if updateTodoReq.ID == 0 || updateTodoReq.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			todo, err := h.svc.UpdateTODO(r.Context(), updateTodoReq.ID, updateTodoReq.Subject, updateTodoReq.Description)

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
			}

			var responseTodo model.UpdateTODOResponse
			responseTodo.TODO = *todo

			jsonError := json.NewEncoder(w).Encode(responseTodo)

			if jsonError != nil {
				fmt.Println(jsonError)
			}

		}
	} else if r.Method == "GET" {
		// クエリパラメータの取得
		prevId := r.URL.Query().Get("prev_id")
		size := r.URL.Query().Get("size")

		var readTodoReq model.ReadTODORequest
		prevIdInt, _ := strconv.Atoi(prevId)
		sizeInt, _ := strconv.Atoi(size)

		prevIdInt64 := int64(prevIdInt)
		sizeInt64 := int64(sizeInt)

		// defaultを与えるのはsizeが無い時じゃなくて0の時？
		if sizeInt64 == 0 {
			sizeInt64 = 5
		}

		readTodoReq.PrevID = prevIdInt64
		readTodoReq.Size = sizeInt64

		todos, _ := h.svc.ReadTODO(r.Context(), readTodoReq.PrevID, readTodoReq.Size)

		// []*model.TODO を []model.TODO に変換する
		var todoSlice []model.TODO
		for _, todo := range todos {
			todoSlice = append(todoSlice, *todo)
		}

		var responseTodo model.ReadTODOResponse
		responseTodo.TODOs = todoSlice

		// sliceがnullだったら空配列を代入する
		if todoSlice == nil {
			responseTodo.TODOs = make([]model.TODO, 0)
		}

		jsonError := json.NewEncoder(w).Encode(responseTodo)

		if jsonError != nil {
			fmt.Println(jsonError)
		}
	} else if r.Method == "DELETE" {
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		var deleteTodoReq model.DeleteTODORequest
		json.Unmarshal(body, &deleteTodoReq)

		if len(deleteTodoReq.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			err := h.svc.DeleteTODO(r.Context(), deleteTodoReq.IDs)

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
			}

			var deleteTodoRes model.DeleteTODOResponse
			jsonError := json.NewEncoder(w).Encode(deleteTodoRes)

			if jsonError != nil {
				fmt.Println(jsonError)
			}
		}

	}
}
