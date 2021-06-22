package tasks

import (
	"log"

	database "github.com/Sandybaby07/graphql-golang-practice/internal/pkg/db/mysql"
	"github.com/Sandybaby07/graphql-golang-practice/internal/users"
)

// #1
type Task struct {
	ID      string
	Title   string
	Content string
	Creater *users.User
	Editor  *users.User
	Status  Status
}

type Status struct {
	Status string
}

// Save new task
func (task Task) Save() int64 {
	// insert new task sql
	stmt, err := database.Db.Prepare("INSERT INTO Tasks(Title,Content,CreaterID,EditorID,Status) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	// Execute insert
	res, err := stmt.Exec(task.Title, task.Content, task.Creater.ID, task.Editor.ID, task.Status.Status)
	if err != nil {
		log.Fatal(err)
	}
	// Result
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	log.Print("Row inserted!")
	return id
}

// Get all tasks
func GetAll() []Task {
	// Query all tasks sql
	stmt, err := database.Db.Prepare("select T.id, T.title, T.Content, T.CreaterID, T.EditorID, T.Status, U.Username from Tasks T inner join Users U on T.CreaterID = U.ID")
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
	var createrId string
	var editorId string
	var status string
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Content, &createrId, &editorId, &status, &username)
		if err != nil {
			log.Fatal(err)
		}
		task.Creater = &users.User{
			ID:       createrId,
			Username: username,
			Password: "",
		}
		task.Editor = &users.User{
			ID:       editorId,
			Username: username,
			Password: "",
		}
		task.Status = Status{
			Status: status,
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return tasks
}

// Delete
func (task Task) Delete() int64 {
	// delete sql DELETE FROM Tasks WHERE `ID`=?
	stmt, err := database.Db.Prepare("DELETE FROM Tasks WHERE `ID` = ? and `CreaterID` = ?")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(task.ID)
	// Execute delete
	res, err := stmt.Exec(task.ID, task.Creater.ID)
	if err != nil {
		log.Fatal(err)
	}
	// Result
	id, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	if id == 0 {
		log.Print("No data!")
	} else {
		log.Print("Row deleted!")
	}
	// Return deleted row quentity
	return id
}

// Modify
func (task Task) Modify() int64 {
	// Update sql
	stmt, err := database.Db.Prepare("UPDATE Tasks SET `Title` = ?,`content` = ?, `EditorID` = ?, `Status` = ? WHERE `ID` = ?")
	if err != nil {
		log.Fatal(err)
	}
	// Execute update
	res, err := stmt.Exec(task.Title, task.Content, task.Editor.ID, task.Status.Status, task.ID)
	if err != nil {
		log.Fatal(err)
	}
	// Result
	id, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	if id == 0 {
		log.Print("No data to update!")
	} else {
		log.Print("Row modify")
	}
	// Return modified row quentity
	return id
}
