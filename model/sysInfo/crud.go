package sysInfo

import (
	"context"

	"github.com/wayne011872/goSterna/model/mgom"
	"github.com/wayne011872/goSterna/log"
	"github.com/wayne011872/systemMonitorServer/dao"
	"github.com/wayne011872/systemMonitorServer/dao/mon"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewCRUD(ctx context.Context, db *mongo.Database, l log.Logger) CRUD {
	return &mongoCRUD{
		mgo: mgom.NewMgoModel(ctx, db, l),
		log:l,
	}
}

type CRUD interface {
	Save(primitive.ObjectID,*dao.SysInfo) error
}

type mongoCRUD struct {
	mgo mgom.MgoDBModel
	log log.Logger
}

func (m *mongoCRUD) Save(orgid primitive.ObjectID,s *dao.SysInfo) error{
	_, err := m.mgo.Save(&mon.SysInfo{
		ID:    primitive.NewObjectID(),
		SysInfo:  s,
	}, nil)
	return err
}