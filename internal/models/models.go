package models

import "time"

var (
	// StatusInProcess указываеь таски, над которыми сейчас работает воркер
	StatusInProcess = "calculating"

	// StatusNew указывает, что объект только что был создан
	StatusNew = "created"

	// StatusFailed указывает, что выражение не решено. Причиной может быть его некорректность
	StatusFailed = "failed"

	// StatusResolved указывает в БД, что результат выражения подсчитан успешно
	StatusResolved = "done"

	// StatusWait указывает на те выражения в БД, результат которых еще не подсчитан
	StatusWait = "awaiting processing"
)

// ContextKey - тип для ключа контекста
type ContextKey string

const (
	// UserIDKey - типизированная константа ID пользователя
	UserIDKey ContextKey = "userID"
)

// ErrorResponse - структура ответа, возвращаемого при ошибке вычислений
type ErrorResponse struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"error_message"`
}

// Expression - структура математического выражения
type Expression struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	Expression   string  `json:"expression"`
	Status       string  `json:"status"`
	Result       float64 `json:"result"`
	ErrorMessage string  `json:"error_message"`
}

// ExpressionRepository — интерфейс для работы с задачами и выражениями
type ExpressionRepository interface {
	GetTask() (*Task, int, error)
	UpdateTaskResult(id int, result float64, err string) error
	GetExpression(id int) (*Expression, error)
	Insert(expr string, userID int) (int, error)
	UpdateStatus(id int, status string)
	UpdateTaskStatus(id int, status string)
}

// RegisteredExpression - структура ответа, возвращаемого при регистрации выражения в оркестраторе
type RegisteredExpression struct {
	ID int `json:"id"`
}

// Request - структура запроса
type Request struct {
	Expression string `json:"expression"`
}

// Response - струтура ответа после успешного завершения программы
type Response struct {
	Expression Expression `json:"expression"`
}

// Task описывает задачу для выполнения
type Task struct {
	ID           int     `json:"ID"`
	ExpressionID int     `json:"ExpressionID"`
	Arg1         float64 `json:"Arg1"`
	Arg2         float64 `json:"Arg2"`
	PrevTaskID1  int     `json:"PrevTaskID1"`
	PrevTaskID2  int     `json:"PrevTaskID2"`
	Operation    string  `json:"Operation"`
	Status       string  `json:"Status"`
	Result       float64 `json:"Result"`
}

// TaskResponse - структура, содержащая одну таску
type TaskResponse struct {
	Task Task `json:"task"`
}

// Token - структура токена, на которые разбивается исходное выражение
type Token struct {
	Value    string
	IsNumber bool
}

// User описывает структуру пользователя
type User struct {
	ID           int       `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}
