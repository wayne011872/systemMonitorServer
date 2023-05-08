package sysInfo

import (
	"context"

	"github.com/wayne011872/goSterna/model/mgom"
	"github.com/wayne011872/goSterna/log"
	"github.com/wayne011872/systemMonitorServer/dao"
	"github.com/wayne011872/systemMonitorServer/dao/mon"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewCRUD(ctx context.Context, db *mongo.Database, l log.Logger) CRUD {
	return &mongoCRUD{
		mgo: mgom.NewMgoModel(ctx, db, l),
		log:l,
	}
}

type CRUD interface {
	Save(*dao.SysInfo) (*mon.SysInfo,error)
	List() []*mon.SysInfo
	SearchList(q bson.M) []*mon.SysInfo
}

type mongoCRUD struct {
	mgo mgom.MgoDBModel
	log log.Logger
}

func (m *mongoCRUD) Save(s *dao.SysInfo) (*mon.SysInfo,error){
	o := & mon.SysInfo{
		ID: primitive.NewObjectID(),
		SysInfo: s,
	}
	_, err := m.mgo.Save(o, nil)
	return o,err
}

func (m *mongoCRUD) List() []*mon.SysInfo {
	result, err := m.mgo.Find(&mon.SysInfo{},bson.M{})
	if err != nil {
		m.log.Warn(err.Error())
		return nil
	}
	return result.([]*mon.SysInfo)
}

func (m *mongoCRUD) SearchList(q bson.M) []*mon.SysInfo {
	result,err:=m.mgo.Find(&mon.SysInfo{},q)
	if err != nil {
		m.log.Warn(err.Error())
		return nil
	}
	return result.([]*mon.SysInfo)
}