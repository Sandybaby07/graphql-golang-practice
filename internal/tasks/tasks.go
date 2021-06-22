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
	Status  Status
}

type Status struct {
	Status string
}

//#2
func (task Task) Save() int64 {
	//#3
	stmt, err := database.Db.Prepare("INSERT INTO Tasks(Title,Content,CreaterID,EditorID,Status) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	//#4
	res, err := stmt.Exec(task.Title, task.Content, task.Creater.ID, task.Editor.ID, task.Status.Status)
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

// delete
func (task Task) Delete() int64 {
	// delete sql DELETE FROM Tasks WHERE `ID`=?
	stmt, err := database.Db.Prepare("DELETE FROM Tasks WHERE `ID` = ? and `CreaterID` = ?")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(task.ID)
	//#4
	res, err := stmt.Exec(task.ID, task.Creater.ID)
	if err != nil {
		log.Fatal(err)
	}
	//#5
	id, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	if id == 0 {
		log.Print("No data!")
	} else {
		log.Print("Row deleted!")
	}
	return id
}

// modify
func (task Task) Modify() int64 {
	// update sql
	stmt, err := database.Db.Prepare("UPDATE Tasks SET `Title` = ?,`content` = ?, `EditorID` = ?, `Status` = ? WHERE `ID` = ?")
	if err != nil {
		log.Fatal(err)
	}
	//#4
	res, err := stmt.Exec(task.Title, task.Content, task.Editor.ID, task.Status.Status, task.ID)
	if err != nil {
		log.Fatal(err)
	}
	//#5
	id, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	if id == 0 {
		log.Print("No data to update!")
	} else {
		log.Print("Row modify")
	}
	return id
}
