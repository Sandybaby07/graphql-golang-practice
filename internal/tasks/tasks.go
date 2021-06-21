package tasks

import (
	"log"

	database "github.com/glyphack/graphlq-golang/internal/pkg/db/mysql"
	"github.com/glyphack/graphlq-golang/internal/users"
)

// #1
type Task struct {
	ID      string
	Title   string
	Content string
	Creater *users.User
	Editor  *users.User
}

//#2
func (task Task) Save() int64 {
	//#3
	stmt, err := database.Db.Prepare("INSERT INTO Tasks(Title,Content,CreaterID,EditorID) VALUES(?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	//#4
	res, err := stmt.Exec(task.Title, task.Content, task.Creater.ID, task.Editor.ID)
	if err != nil {
		log.Fatal(err)
	}
	//#5
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	log.Print("Row inserted!")
	return id
}

func GetAll() []Task {
	stmt, err := database.Db.Prepare("select T.id, T.title, T.Content, T.CreaterID, T.EditorID, U.Username from Tasks T inner join Users U on T.CreaterID = U.ID")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var tasks []Task
	var username string
	var id string
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Content, &id, &username)
		if err != nil {
			log.Fatal(err)
		}
		task.Creater = &users.User{
			ID:       id,
			Username: username,
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return tasks
}
