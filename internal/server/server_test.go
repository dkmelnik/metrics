package server

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestServer_Run(t *testing.T) {
	addr := "127.0.0.1:8080"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := NewServer(addr, handler)

	go func() {
		err := server.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("server.Run() returned an unexpected error: %v", err)
		}
	}()

	// Ждем некоторое время, чтобы сервер успел запуститься
	time.Sleep(time.Second)

	// Отправляем тестовый запрос на сервер
	req, err := http.NewRequest("GET", "http://"+addr+"/", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем код ответа сервера
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", resp.StatusCode, http.StatusOK)
	}

	// Останавливаем сервер
	if err := server.app.Shutdown(context.TODO()); err != nil {
		t.Errorf("failed to shutdown server: %v", err)
	}
}
