package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

type Request struct {
	Expression string `json:"expression"`
}

type Id struct {
	Id int `json:"id"`
}

type Orchestrator struct {
	exprList          []Expression
	taskList          []Task
	taskQueue         []*Task
	mu                sync.Mutex
	counterExpression int
	counterTask       int
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		exprList:  []Expression{},
		taskList:  []Task{},
		taskQueue: []*Task{},
	}
}

type Expression struct {
	Id     int      `json:"id"`
	Exp    string   `json:"expression"`
	Status string   `json:"status"`
	Result float64  `json:"result"`
	AST    *ASTNode `json:"-"`
}

type Task struct {
	ID            string   `json:"id"`
	ExprID        string   `json:"-"`
	Arg1          float64  `json:"arg1"`
	Arg2          float64  `json:"arg2"`
	Operation     string   `json:"operation"`
	OperationTime int      `json:"operation_time"`
	Node          *ASTNode `json:"-"`
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println(err, "fssfd")
		log.Fatal("Ошибка загрузки .env файла")
	}
}

func Valid(e string) bool {
	valid_chars := "1234567890+-*/()"
	// чек на посторонние символы и равное кол-во открывающихся и закрывающихся скобок
	c1 := 0
	c2 := 0
	for i := range len(e) {
		if !strings.ContainsRune(valid_chars, rune(e[i])) {
			log.Printf("Невалидные символы")
			return false
		}
		if string(e[i]) == "(" {
			c1 += 1
		} else if string(e[i]) == ")" {
			c2 += 1
		}
	}
	if c1 != c2 {
		log.Printf("Неравное кол-во скобок")
		return false
	}

	// чек на неправильную расстановку
	for i := range len(e) - 1 {
		if (string(e[i]) == "+" || string(e[i]) == "-" || string(e[i]) == "/" || string(e[i]) == "*" || string(e[i]) == "(" || string(e[i]) == ")") && (string(e[i+1]) == "+" || string(e[i+1]) == "*" || string(e[i+1]) == "/" || string(e[i+1]) == "-" || string(e[i+1]) == ")") {
			log.Printf("Невалидные знаки")
			return false
		}
	}
	//чек ласт символ
	if string(e[len(e)-1]) == "+" || string(e[len(e)-1]) == "-" || string(e[len(e)-1]) == "*" || string(e[len(e)-1]) == "/" || string(e[len(e)-1]) == "(" {
		log.Printf("Неверный последний символ")
	}
	return true
}

func (o *Orchestrator) createHandler(w http.ResponseWriter, r *http.Request) {
	req := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)
	req.Expression = strings.Replace(req.Expression, " ", "", -1)
	if err != nil || !Valid(req.Expression) {
		http.Error(w, "Невалидные данные", http.StatusUnprocessableEntity)
		return
	}
	ast, _ := ParseAST(req.Expression)

	o.mu.Lock()
	defer o.mu.Unlock()
	o.counterExpression++
	expr := Expression{
		Id:     o.counterExpression,
		Exp:    req.Expression,
		Status: "cooking",
		AST:    ast,
	}
	o.exprList = append(o.exprList, expr)
	o.AstParseExpression(expr)

	w.WriteHeader(http.StatusCreated)
	log.Printf("Выражение принято на вычисление")
	json.NewEncoder(w).Encode(Id{Id: expr.Id})
}

func (o *Orchestrator) getTaskHandler(w http.ResponseWriter, _ *http.Request) {

	o.mu.Lock()
	defer o.mu.Unlock()

	if len(o.taskQueue) == 0 {
		http.Error(w, `{"error":"таски закончились"}`, http.StatusNotFound)
		return
	}

	task := o.dequeueTask()
	o.updateGetExpressionStatus(task.ExprID, "in_progress")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
}

func (o *Orchestrator) dequeueTask() *Task {
	task := o.taskQueue[0]
	o.taskQueue = o.taskQueue[1:]
	return task
}

func (o *Orchestrator) updateGetExpressionStatus(exprID string, status string) {
	for i := range o.exprList {
		id, _ := strconv.Atoi(exprID)
		if o.exprList[i].Id == id {
			o.exprList[i].Status = status
			break
		}
	}
}

func (o *Orchestrator) postTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req Response
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"неваилдные данные"}`, http.StatusUnprocessableEntity)
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	task, taskIndex := o.findTaskByID(strconv.Itoa(req.Id))
	if task == nil {
		http.Error(w, `{"error":"таск не найден"}`, http.StatusNotFound)
		return
	}
	o.updateASTNode(task.Node, req.Res)
	o.taskList = append(o.taskList[:taskIndex], o.taskList[taskIndex+1:]...)
	exprIndex := o.findExpressionIndexByID(task.ExprID)
	if exprIndex != -1 {
		o.AstParseExpression(o.exprList[exprIndex])
		o.updatePostExpressionStatus(o.exprList[exprIndex])
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"резльтат принят"}`))
}

func (o *Orchestrator) updateASTNode(node *ASTNode, result float64) {
	node.IsLeaf = true
	node.Value = result
}

func (o *Orchestrator) findTaskByID(taskID string) (*Task, int) {
	for i := range o.taskList {
		if o.taskList[i].ID == taskID {
			return &o.taskList[i], i
		}
	}
	return nil, -1
}

func (o *Orchestrator) findExpressionIndexByID(exprID string) int {
	for i := range o.exprList {
		a, _ := strconv.Atoi(exprID)
		if o.exprList[i].Id == a {
			return i
		}
	}
	return -1
}

func (o *Orchestrator) updatePostExpressionStatus(expr Expression) {
	if expr.AST.IsLeaf {
		expr.Status = "completed"
		expr.Result = expr.AST.Value
	}
}

func (o *Orchestrator) AstParseExpression(expr Expression) {
	//рекурсивный обход AST
	var bypass func(node *ASTNode)
	bypass = func(node *ASTNode) {
		if node == nil || node.IsLeaf {
			return
		}

		bypass(node.Right)
		bypass(node.Left)

		if node.Left != nil && node.Right != nil && node.Left.IsLeaf && node.Right.IsLeaf {
			id := strconv.Itoa(expr.Id)
			task := o.createTask(node, id)

			o.taskList = append(o.taskList, *task)
			o.taskQueue = append(o.taskQueue, task)
		}
	}

	// начинаем обход с корневого узла AST
	bypass(expr.AST)
}

func (o *Orchestrator) createTask(node *ASTNode, exprID string) *Task {
	o.counterTask++
	taskID := fmt.Sprintf("%d", o.counterTask)
	opTime := o.getOperationTime(node.Operator)
	return &Task{
		ID:            taskID,
		ExprID:        exprID,
		Arg1:          node.Left.Value,
		Arg2:          node.Right.Value,
		Operation:     node.Operator,
		OperationTime: opTime,
		Node:          node,
	}
}
func (o *Orchestrator) getOperationTime(operator string) int {
	var envVar string
	switch operator {
	case "+":
		envVar = "TIME_ADDITION_MS"
	case "-":
		envVar = "TIME_SUBTRACTION_MS"
	case "*":
		envVar = "TIME_MULTIPLICATIONS_MS"
	case "/":
		envVar = "TIME_DIVISIONS_MS"
	default:
		return 1000
	}
	timeStr := os.Getenv(envVar)
	if timeStr == "" {
		log.Printf("Переменная %s не сущесвтует", envVar)
		return 1000
	}

	time, err := strconv.Atoi(timeStr)
	if err != nil {
		log.Printf("Невалидные данные %s: %s", envVar, timeStr)
		return 1000
	}

	return time
}
func (o *Orchestrator) expressionsHandler(w http.ResponseWriter, r *http.Request) {
	o.mu.Lock()
	defer o.mu.Unlock()
	exprs := make([]Expression, 0, len(o.exprList))
	for _, expr := range o.exprList {
		if expr.AST != nil && expr.AST.IsLeaf {
			expr.Status = "completed"
			expr.Result = expr.AST.Value
		}
		exprs = append(exprs, expr)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": exprs})
}

func (o *Orchestrator) expressionIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	if id == "" {
		http.Error(w, `{"error":"неверный айди"}`, http.StatusBadRequest)
		return
	}
	idi, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error":"выражение не существует"}`, http.StatusNotFound)
		return
	}
	o.mu.Lock()
	if idi > len(o.exprList) {
		http.Error(w, `{"error":"выражение не существует"}`, http.StatusNotFound)
		return
	}
	expr := o.exprList[idi-1]
	o.mu.Unlock()

	if expr.AST != nil && expr.AST.IsLeaf {
		expr.Status = "completed"
		expr.Result = expr.AST.Value
	}

	response := struct {
		Expression Expression `json:"expression"`
	}{
		Expression: expr,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error":"ошибка сервера"}`, http.StatusInternalServerError)
	}
}

func (o *Orchestrator) Run() error {
	http.HandleFunc("/api/v1/calculate", o.createHandler)
	http.HandleFunc("/api/v1/expressions", o.expressionsHandler)
	http.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			o.getTaskHandler(w, r)
		}
		if r.Method == http.MethodPost {
			o.postTaskHandler(w, r)
		}
	})
	http.HandleFunc("/api/v1/expressions/", o.expressionIDHandler)
	log.Printf("Сервер запущен")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
		return nil
	} else {
		return err
	}

}
