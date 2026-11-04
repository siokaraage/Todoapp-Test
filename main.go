package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

const fileName = "todos.json"

// ファイルから読み込む。無ければ空の一覧を返す
func load() ([]Todo, error) {
	data, err := os.ReadFile(fileName)
	if os.IsNotExist(err) {
		return []Todo{}, nil
	}
	if err != nil {
		return nil, err
	}
	var todos []Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

// ファイルへ保存する
func save(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func main() {
	if len(os.Args[1]) < 2 {
	fmt.Println("使い方: add(内容) | list | done (ID)")
	os.Exit(1)
}

todos, err := load()
if err != nil {
	fmt.Println("読み込みに失敗しました。", err)
	os.Exit(1)
}

switch os.Args[1] {
case "add":
	title := strings.Join(os.Args[2:], "")
	if title == "" {
		fmt.Println("内容を入力してください")
		os.Exit(1)
	}

	maxID := 0
	for _,t := range todos {
		if t.ID > maxID {
		maxID = t.ID
		}
	}
	todos = append(todos, Todo{ID: maxID + 1, Title: title})
	if err := save(todos); err != nil {
		fmt.Println("保存に失敗しました。", err)
		os.Exit(1)
	}
	fmt.Println("追加しました:", title)

case "list":
	for _, t := range todos {
		mark := "[]"
		if t.Done {
			mark = "[x]"
		}
		fmt.Printf("%d %s %s\n", t.ID, mark, t.Title)
	}

case "done":
	if len(os.Args) < 3 {
		fmt.Println("IDを指定してください")
		os.Exit(1)
	}
id, err := strconv.Atoi(os.Args[2])
if err != nil {
	fmt.Println("IDは数字で指定してください")
	os.Exit(1)
}
found := false
for i := range todos {
	if todos[i].ID == id {
		todos[i].Done = true
		found = true
	}
}
if !found {
	fmt.Println("そのIDは見つかりません")
	os.Exit(1)
}
if err := save(todos); err != nil {
	fmt.Println("保存に失敗:", err)
	os.Exit(1)
}
fmt.Println("完了にしました:", id)

default:
	fmt.Println("不明なコマンドです:", os.Args[1])
	os.Exit(1)
}
}