package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mizhidajiangyou/go-linux/conf"
	"github.com/mizhidajiangyou/go-linux/log"
	"os"
	"reflect"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

//MySQL：username:password@tcp(hostname:port)/databasename
//
//PostgreSQL：user=username password=password host=hostname port=port dbname=databasename
//
//SQLite：file:/path/to/database.db
//
//Microsoft SQL Server：sqlserver://username:password@hostname:port?database=databasename
//
//Oracle：username/password@ip:port/service_name

var databaseConf *DatabaseConf
var dsn string

//var Engine *xorm.Engine

type DatabaseConf struct {
	Type     string `json:"type" env:"DB_TYPE"`
	Host     string `json:"host" env:"DB_HOST"`
	Port     string `json:"port" env:"DB_PORT"`
	User     string `json:"user" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASSWD"`
	DBName   string `json:"name" env:"DB_NAME"`
	DBFile   string `json:"db_file" env:"DB_FILE"`
}

//从环境变量中获取
func initConfByEnv() {
	cf := &DatabaseConf{}

	// 获取结构体类型
	confType := reflect.TypeOf(cf).Elem()

	// 遍历结构体字段
	for i := 0; i < confType.NumField(); i++ {
		field := confType.Field(i)
		envTag := field.Tag.Get("env")

		// 从环境变量获取值
		value := os.Getenv(envTag)
		if value != "" {
			// 赋值给结构体字段
			fieldValue := reflect.ValueOf(cf).Elem().FieldByName(field.Name)
			if fieldValue.IsValid() && fieldValue.CanSet() {
				fieldValue.SetString(value)
			}
		}

	}

	databaseConf = cf
}

//从配置文件中获取
func initConfByConf(filePath string) {

	// 获取结构体类型
	cf, _ := conf.ReadYaml(filePath, DatabaseConf{})

	structValue := reflect.ValueOf(databaseConf).Elem()
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)
		if value, ok := cf["DatabaseConf"].(map[string]interface{})[field.Name]; ok {
			if reflect.TypeOf(value).AssignableTo(field.Type) {
				fieldValue.Set(reflect.ValueOf(value))
				log.Debugf(fmt.Sprintf("set value: %s ", value))
			}
		}
	}

}

func ReadDBConf(filePath string) {
	initConfByConf(filePath)
}

func StartDb() (engine *xorm.Engine, err error) {
	// 设置数据库连接字符串

	switch databaseConf.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4",
			databaseConf.User,
			databaseConf.Password,
			databaseConf.Host,
			databaseConf.Port,
			databaseConf.DBName,
		)
	default:
		log.Warnf(fmt.Sprintf("未识别类型 %s ，将使用默认类型mysql.", databaseConf.Type))
		databaseConf.Type = "mysql"
	}
	// 创建引擎
	engine, err = xorm.NewEngine(databaseConf.Type, dsn)
	if err != nil {
		log.Fatal(err)
		return
	}

	// 设置表名规则
	engine.SetMapper(names.SameMapper{})

	// 测试连接
	if err = engine.Ping(); err != nil {
		log.Fatal(err)
		return
	}

	log.Infof("Connected to MySQL!")
	return

}

func init() {
	initConfByEnv()
}
