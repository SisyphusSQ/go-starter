package mongo

import (
	"context"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"

	"go-starter/config"
	"go-starter/internal/lib/mongodb"
	"go-starter/internal/models/do/mongo/task"
)

type TaskListRepository interface {
	GetAll(ctx context.Context) (tasks []*task.TaskList, err error)
	GetByName(ctx context.Context, name string) (task task.TaskList, err error)
	GetByCondition(ctx context.Context, cond bson.M) (task []*task.TaskList, err error)

	Create(ctx context.Context, task task.TaskList) (err error)
	CreateBatch(ctx context.Context, tasks []*task.TaskList) (err error)

	DeleteById(ctx context.Context, task task.TaskList) (err error)
}

type mongoTaskListRepo struct {
	client *qmgo.QmgoClient
}

func NewTaskListRepository(c config.MongoDB, coll string) TaskListRepository {
	cli, err := mongodb.New(c, coll)
	if err != nil {
		panic(err)
	}

	return &mongoTaskListRepo{client: cli}
}

func (m *mongoTaskListRepo) GetAll(ctx context.Context) (tasks []*task.TaskList, err error) {
	err = m.client.Find(ctx, bson.M{}).All(&tasks)
	return
}

func (m *mongoTaskListRepo) GetByName(ctx context.Context, name string) (task task.TaskList, err error) {
	err = m.client.Find(ctx, bson.M{"name": name}).One(&task)
	return
}

func (m *mongoTaskListRepo) GetByCondition(ctx context.Context, cond bson.M) (task []*task.TaskList, err error) {
	err = m.client.Find(ctx, cond).All(&task)
	return
}

func (m *mongoTaskListRepo) Create(ctx context.Context, task task.TaskList) (err error) {
	_, err = m.client.InsertOne(ctx, &task)
	return
}

func (m *mongoTaskListRepo) CreateBatch(ctx context.Context, tasks []*task.TaskList) (err error) {
	_, err = m.client.InsertMany(ctx, &tasks)
	return
}

func (m *mongoTaskListRepo) UpdateByName(ctx context.Context, name string, updates bson.M) (err error) {
	condition := bson.M{
		"name": name,
		"isDelete": bson.M{
			"$ne": true,
		},
	}

	err = m.client.UpdateOne(ctx, condition, bson.M{"$set": updates})
	return
}

func (m *mongoTaskListRepo) DeleteById(ctx context.Context, task task.TaskList) (err error) {
	//err = m.client.RemoveId(ctx, task.Id)

	del := bson.M{"isDelete": true}
	err = m.client.UpdateId(ctx, task.Id, bson.M{"$set": del})
	return
}
