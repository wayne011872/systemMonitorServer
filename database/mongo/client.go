package mongo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetMgoDBClient() (MongoDBClient,error) {
	mc := &MongoConf{}
	mgoclt, err := mc.NewMongoDBClient()
	if err != nil {
		return nil,err
	}
	return mgoclt,nil
}

type MongoDI interface {
	NewMongoDBClient()(MongoDBClient, error)
	SetAuthCredential()
	GetConfig()
	GetUri()string
}

type MongoConf struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
	ErrDataBase string `yaml:"errdatabase"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`

	authCredential  options.Credential
}

func (mc *MongoConf) GetConfig() {
	configPath := os.Getenv(("CONFIG_PATH"))
	configName := os.Getenv(("CONFIG_NAME"))
	configType := os.Getenv(("CONFIG_TYPE"))
	if configPath == "" {
		panic("沒有設定CONFIG_PATH參數")
	}
	if configName == "" {
		panic("沒有設定CONFIG_NAME參數")
	}
	if configType == "" {
		panic("沒有設定CONFIG_TYPE參數")
	}
	vip := viper.New()
	vip.AddConfigPath(configPath)
	vip.SetConfigName(configName)
	vip.SetConfigType(configType)
	if err := vip.ReadInConfig(); err != nil {
		panic(err)
	}
	err := vip.UnmarshalKey("mongo", &mc)
	if err != nil {
		panic(err)
	}
}
func (mc *MongoConf) SetAuthCredential() {
	if mc.User != "" && mc.Pass != "" {
		mc.authCredential = options.Credential{
			AuthMechanism: "SCRAM-SHA-256",
			AuthSource: "admin",
			Username: mc.User,
			Password: mc.Pass,
		}
	}
}

func (mc *MongoConf) GetUri() string {
	return mc.Uri
}

func (mc *MongoConf) NewMongoDBClient() (MongoDBClient, error) {
	mc.GetConfig()
	if mc.Uri == "" {
		panic("mongo uri is not set!")
	}
	if mc.Database == "" {
		panic("mongo DataBase is not set!")
	}
	var clientOpts *options.ClientOptions
	if mc.User != "" && mc.Pass != "" {
		mc.SetAuthCredential()
		clientOpts = options.Client().ApplyURI(mc.GetUri()).SetAuth(mc.authCredential)
	}else{
		clientOpts = options.Client().ApplyURI(mc.GetUri())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		cancel()
		return nil,err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		cancel()
		return nil, err
	}
	db := client.Database(mc.Database)
	errDb := client.Database(mc.ErrDataBase)
	return &mgoClientImpl{
		clt:    client,
		ctx:    ctx,
		cancel: cancel,
		db:     db,
		errDb:  errDb,
	}, nil
}

type mgoClientImpl struct {
	clt     *mongo.Client
	ctx     context.Context
	cancel  context.CancelFunc
	session mongo.Session
	db      *mongo.Database
	errDb   *mongo.Database
}

func (m *mgoClientImpl) WithSession(f func(sc mongo.SessionContext) error) error {
	if m.session != nil {
		return nil
	}
	session, err := m.clt.StartSession()
	if err != nil {
		m.cancel()
		return err
	}
	if err := session.StartTransaction(); err != nil {
		return err
	}
	m.session = session
	return mongo.WithSession(m.ctx, m.session, f)
}

func (m *mgoClientImpl) Close() {
	if m == nil {
		return
	}
	if m.session != nil {
		m.session.EndSession(m.ctx)
	}
	if m.clt != nil {
		err := m.clt.Disconnect(m.ctx)
		if err != nil {
			fmt.Println("disconnect error: " + err.Error())
		}
	}
	m.cancel()
}

func (m *mgoClientImpl) AbortTransaction(sc mongo.SessionContext) error {
	return m.session.AbortTransaction(sc)
}

func (m *mgoClientImpl) CommitTransaction(sc mongo.SessionContext) error {
	return m.session.CommitTransaction(sc)
}
func (m *mgoClientImpl) GetDB() *mongo.Database {
	return m.db
}
func (m *mgoClientImpl) GetErrDB() *mongo.Database {
	return m.errDb
}

func (m *mgoClientImpl) Ping() error {
	return m.clt.Ping(m.ctx, readpref.Primary())
}

func (m *mgoClientImpl) GetCtx() context.Context {
	return m.ctx
}

type MongoDBClient interface {
	WithSession(f func(sc mongo.SessionContext) error) error
	AbortTransaction(sc mongo.SessionContext) error
	CommitTransaction(sc mongo.SessionContext) error
	Close()
	GetDB() *mongo.Database
	GetErrDB() *mongo.Database
	Ping() error
	GetCtx() context.Context
}
