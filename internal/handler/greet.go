package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/your-username/greeting-service/internal/store"
)

// Профессии для выпадающего списка
var professions = []string{
	"Инженер",
	"Врач",
	"Учитель",
	"Программист",
	"Дизайнер",
	"Менеджер",
	"Другое",
}

// Данные для главного шаблона (форма + таблица)
type pageData struct {
	Professions []string
	Records     []store.Record
}

// FormHandler — главная страница с формой и таблицей
func FormHandler(w http.ResponseWriter, r *http.Request) {
	data := pageData{
		Professions: professions,
		Records:     store.GetAll(),
	}

	tmpl := `<!DOCTYPE html>
<html>
<head><title>Приветствие</title></head>
<body>
    <h2>Введите данные</h2>
    <form action="/greet" method="POST">
        <input type="text" name="first_name" placeholder="Имя" required>
        <input type="text" name="last_name" placeholder="Фамилия" required>
        <select name="profession" required>
            <option value="">Выберите профессию</option>
            {{range .Professions}}
            <option value="{{.}}">{{.}}</option>
            {{end}}
        </select>
        <button type="submit">Поздороваться</button>
    </form>
    
    <hr>
    <h2>Все записи</h2>
    {{if .Records}}
    <table border="1">
        <tr>
            <th>Имя</th>
            <th>Фамилия</th>
            <th>Профессия</th>
        </tr>
        {{range .Records}}
        <tr>
            <td>{{.FirstName}}</td>
            <td>{{.LastName}}</td>
            <td>{{.Profession}}</td>
        </tr>
        {{end}}
    </table>
    {{else}}
    <p>Пока нет записей.</p>
    {{end}}
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.Must(template.New("main").Parse(tmpl))
	t.Execute(w, data)
}

// GreetFormHandler обрабатывает POST от HTML-формы и перенаправляет на главную
func GreetFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	profession := r.FormValue("profession")

	if firstName == "" || lastName == "" || profession == "" {
		http.Error(w, "Имя, фамилия и профессия обязательны", http.StatusBadRequest)
		return
	}

	// Сохраняем запись
	store.Add(store.Record{
		FirstName:  firstName,
		LastName:   lastName,
		Profession: profession,
	})

	// Перенаправляем на главную, чтобы избежать повторной отправки формы при обновлении страницы
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// JSON структуры
type greetRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Profession string `json:"profession"`
}

type greetResponse struct {
	Message string `json:"message"`
}

// GreetAPIHandler обрабатывает JSON API запросы и сохраняет запись
func GreetAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req greetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.FirstName == "" || req.LastName == "" || req.Profession == "" {
		http.Error(w, "first_name, last_name, profession are required", http.StatusBadRequest)
		return
	}

	store.Add(store.Record{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Profession: req.Profession,
	})

	resp := greetResponse{
		Message: fmt.Sprintf("Привет, %s %s (%s)!", req.FirstName, req.LastName, req.Profession),
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

// ListHTMLHandler — отдельная страница с таблицей (оставлена для обратной совместимости)
func ListHTMLHandler(w http.ResponseWriter, r *http.Request) {
	records := store.GetAll()

	tmpl := `<!DOCTYPE html>
<html>
<head><title>Список записей</title></head>
<body>
    <h2>Все записи</h2>
    <table border="1">
        <tr>
            <th>Имя</th>
            <th>Фамилия</th>
            <th>Профессия</th>
        </tr>
        {{range .}}
        <tr>
            <td>{{.FirstName}}</td>
            <td>{{.LastName}}</td>
            <td>{{.Profession}}</td>
        </tr>
        {{end}}
    </table>
    <p><a href="/">Вернуться к форме</a></p>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.Must(template.New("list").Parse(tmpl))
	t.Execute(w, records)
}

// ListAPIHandler возвращает JSON со списком всех записей
func ListAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	records := store.GetAll()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(records)
}
