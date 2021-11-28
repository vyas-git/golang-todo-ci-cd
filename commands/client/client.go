package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-grpc/internal/rpc"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

const (
	address = "todo-rpc-service:9090"
)

type TodoPageData struct {
	PageTitle string
	Todos     []*rpc.Todo
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(pwd)
	dialRPC()
	initRoutes()
}
func dialRPC() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
}
func initRoutes() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Method("GET", "/todos", handler(indexTodos))
	r.Method("POST", "/addtodo", handler(createTodo))
	r.Method("GET", "/getTodos", handler(getTodos))
	r.Method("DELETE", "/deleteTodo/{id}", handler(deleteTodo))
	log.Printf("localhost listening at %v", ":3000")

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatalf("failed to serve localhost:3000: %v", err)
		return
	}

}
func handler(f func(http.ResponseWriter, *http.Request, rpc.TodoServiceClient) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		c := rpc.NewTodoServiceClient(conn)
		if err := f(w, r, c); err != nil {
			fmt.Println("Error", err)
		}
	})
}
func createTodo(w http.ResponseWriter, r *http.Request, c rpc.TodoServiceClient) error {
	ctx := context.Background()
	fmt.Println(r.Body)
	// define custom type
	var newTodo rpc.Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Decode error! please check your JSON formating.")
		return err
	}
	result, err := c.Create(ctx, &rpc.NewTodo{
		Title:   newTodo.Title,
		Content: &newTodo.Title,
	})
	if err != nil {
		log.Fatalf("could not create todo: %v", err)
	}
	log.Printf("todo detail: %v", result)
	return nil
}

func indexTodos(w http.ResponseWriter, r *http.Request, c rpc.TodoServiceClient) error {
	w.Header().Set("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("todolist.html"))

	ctx := context.Background()
	result, err := c.Index(ctx, &rpc.Empty{})
	if err != nil {
		log.Fatalf("could not index todo")
	}
	data := TodoPageData{
		PageTitle: "My TODOs list",
		Todos:     result.Items,
	}
	tmpl.Execute(w, data)

	//w.Write([]byte("<h1>hii</h1>"))
	//	json.NewEncoder(w).Encode(result)
	log.Println("list todos")
	log.Println(result)
	return nil
}
func getTodos(w http.ResponseWriter, r *http.Request, c rpc.TodoServiceClient) error {
	ctx := context.Background()
	result, err := c.Index(ctx, &rpc.Empty{})
	if err != nil {
		log.Fatalf("could not index todo")
		return err
	}
	json.NewEncoder(w).Encode(result)
	return nil
}
func deleteTodo(w http.ResponseWriter, r *http.Request, c rpc.TodoServiceClient) error {
	ctx := context.Background()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Fatalf("could not delete todo")
		return err
	}
	_, err = c.Delete(ctx, &rpc.TodoId{
		Id: int32(id),
	})
	if err != nil {
		log.Fatalf("could not delete todo")
		return err
	}
	result, err := c.Index(ctx, &rpc.Empty{})
	if err != nil {
		log.Fatalf("could not index todo")
		return err
	}
	json.NewEncoder(w).Encode(result)
	return nil
}

func showTodo(c rpc.TodoServiceClient, ctx context.Context) {
	result, err := c.Show(ctx, &rpc.TodoId{
		Id: 2,
	})
	if err != nil {
		log.Fatalf("could not delete todo")
	}

	fmt.Println(result)
}
