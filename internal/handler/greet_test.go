package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/n1kno/practic8/internal/store"
)

func TestFormHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	FormHandler(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "<form") {
		t.Error("response should contain form")
	}
	if !strings.Contains(body, "Инженер") {
		t.Error("profession options missing")
	}
	if !strings.Contains(body, "Пока нет записей") {
		t.Error("expected 'Пока нет записей' when no records")
	}
}

func TestGreetFormHandlerRedirect(t *testing.T) {
	form := strings.NewReader("first_name=Иван&last_name=Петров&profession=Инженер")
	req := httptest.NewRequest("POST", "/greet", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	GreetFormHandler(w, req)

	// Ожидаем редирект 303
	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303 See Other, got %d", w.Code)
	}
	location := w.Header().Get("Location")
	if location != "/" {
		t.Errorf("expected redirect to /, got %s", location)
	}

	// Проверяем, что запись сохранилась
	all := store.GetAll()
	found := false
	for _, rec := range all {
		if rec.FirstName == "Иван" && rec.LastName == "Петров" && rec.Profession == "Инженер" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("record not saved correctly, records: %+v", all)
	}
}

func TestGreetAPIHandler(t *testing.T) {
	body := `{"first_name":"Мария","last_name":"Иванова","profession":"Врач"}`
	req := httptest.NewRequest("POST", "/api/greet", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	GreetAPIHandler(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"message":"Привет, Мария Иванова (Врач)!"`) {
		t.Errorf("unexpected body: %s", w.Body.String())
	}

	// Проверяем сохранение
	all := store.GetAll()
	found := false
	for _, rec := range all {
		if rec.FirstName == "Мария" && rec.LastName == "Иванова" {
			found = true
			break
		}
	}
	if !found {
		t.Error("record not saved from API")
	}
}

func TestListHTMLHandler(t *testing.T) {
	// Добавим запись для теста
	store.Add(store.Record{FirstName: "Петр", LastName: "Сидоров", Profession: "Менеджер"})
	req := httptest.NewRequest("GET", "/list", nil)
	w := httptest.NewRecorder()
	ListHTMLHandler(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Петр") || !strings.Contains(body, "Менеджер") {
		t.Errorf("table missing expected data: %s", body)
	}
}

func TestListAPIHandler(t *testing.T) {
	store.Add(store.Record{FirstName: "Анна", LastName: "Котова", Profession: "Дизайнер"})
	req := httptest.NewRequest("GET", "/api/list", nil)
	w := httptest.NewRecorder()
	ListAPIHandler(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"first_name":"Анна"`) {
		t.Errorf("unexpected JSON: %s", w.Body.String())
	}
}
