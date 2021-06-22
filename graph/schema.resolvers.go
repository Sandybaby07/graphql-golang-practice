package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/glyphack/graphlq-golang/graph/generated"
	"github.com/glyphack/graphlq-golang/graph/model"
	"github.com/glyphack/graphlq-golang/internal/auth"
	"github.com/glyphack/graphlq-golang/internal/links"
	"github.com/glyphack/graphlq-golang/internal/tasks"
	"github.com/glyphack/graphlq-golang/internal/users"
	"github.com/glyphack/graphlq-golang/pkg/jwt"
)

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Link{}, fmt.Errorf("access denied")
	}
	var link links.Link
	link.Title = input.Title
	link.Address = input.Address
	link.User = user
	linkId := link.Save()
	grahpqlUser := &model.User{
		ID:   user.ID,
		Name: user.Username,
	}
	return &model.Link{ID: strconv.FormatInt(linkId, 10), Title: link.Title, Address: link.Address, User: grahpqlUser}, nil
}

func (r *mutationResolver) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Task{}, fmt.Errorf("access denied")
	}
	var task tasks.Task
	task.Title = input.Title
	task.Content = input.Content
	task.Creater = user
	task.Editor = user
	task.Status.Status = string(model.StatusPending)
	taskId := task.Save()
	grahpqlCreater := &model.User{
		ID:   user.ID,
		Name: user.Username,
	}
	grahpqlEditor := &model.User{
		ID:   user.ID,
		Name: user.Username,
	}
	return &model.Task{ID: strconv.FormatInt(taskId, 10), Title: task.Title, Content: task.Content, Creater: grahpqlCreater, Editor: grahpqlEditor, Status: model.StatusPending}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	user.Role.Role = input.Role.String()
	user.Create()
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) DeleteTask(ctx context.Context, input model.DeleteTask) (string, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return "", fmt.Errorf("access denied")
	}
	var deleteTask tasks.Task
	deleteTask.ID = input.ID
	deleteTask.Creater = user
	delete := deleteTask.Delete()
	return strconv.FormatInt(delete, 10), nil
}

func (r *mutationResolver) ModifyTask(ctx context.Context, input model.ModifyTask) (string, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return "", fmt.Errorf("access denied")
	}
	var modifyTask tasks.Task
	modifyTask.ID = input.ID
	modifyTask.Creater = user
	modifyTask.Title = input.Title
	modifyTask.Content = input.Content
	modifyTask.Editor = user
	modifyTask.Status.Status = input.Status.String()
	modify := modifyTask.Modify()
	return strconv.FormatInt(modify, 10), nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		// 1
		return "", &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	var resultLinks []*model.Link
	var dbLinks []links.Link
	dbLinks = links.GetAll()
	for _, link := range dbLinks {
		grahpqlUser := &model.User{
			Name: link.User.Username,
		}
		resultLinks = append(resultLinks, &model.Link{ID: link.ID, Title: link.Title, Address: link.Address, User: grahpqlUser})
	}
	return resultLinks, nil
}

func (r *queryResolver) Task(ctx context.Context) ([]*model.Task, error) {
	var resultTasks []*model.Task
	var dbTasks []tasks.Task
	dbTasks = tasks.GetAll()
	for _, task := range dbTasks {
		grahpqlCreater := &model.User{
			Name: task.Creater.Username,
		}
		grahpqlEditor := &model.User{
			Name: task.Editor.Username,
		}
		resultTasks = append(resultTasks, &model.Task{ID: task.ID, Title: task.Title, Content: task.Content, Creater: grahpqlCreater, Editor: grahpqlEditor, Status: model.StatusPending})
	}
	return resultTasks, nil
}

func (r *queryResolver) User(ctx context.Context) ([]*model.User, error) {
	var resultUsers []*model.User
	var dbUsers []users.User
	dbUsers = users.GetAll()
	for _, user := range dbUsers {
		resultUsers = append(resultUsers, &model.User{ID: user.ID, Name: user.Username, Role: model.Role(user.Role.Role)})
	}
	return resultUsers, nil
}

func (r *queryResolver) Staff(ctx context.Context) ([]*model.User, error) {
	var resultUsers []*model.User
	var dbStaff []users.User
	dbStaff = users.GetStaff()
	for _, user := range dbStaff {
		resultUsers = append(resultUsers, &model.User{ID: user.ID, Name: user.Username, Role: model.Role(user.Role.Role)})
	}
	return resultUsers, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
