package orm

import (
	"films/utils"
	gormPlugin "films/utils/go2sky_gorm"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-micro.dev/v4/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"time"
)

var conn *gorm.DB
var conf config.Config

type Json_return struct {
	Code int
	Msg  string
	Data any
}

func init() {
	conf = utils.Instance()
}

//各个服务初始化调用的数据库连接
func Initconf(name, service_name string) *gorm.DB {

	if conf.Get("hosts", name, "type").String("") == "mysql" {
		user := conf.Get("hosts", name, "user").String("")
		password := conf.Get("hosts", name, "password").String("")
		host := conf.Get("hosts", name, "address").String("")
		port := conf.Get("hosts", name, "port").Int(3306)
		dbname := conf.Get("hosts", name, "dbname").String("")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, strconv.Itoa(port), dbname)

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             5 * time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Warn,     // 日志级别
				IgnoreRecordNotFoundError: false,           // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  true,            // 禁用彩色打印
			},
		)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})

		conn = db

		if err != nil {
			panic(err)
		}
		fmt.Println("引入连接数据库")
		sqlconn, _ := db.DB()
		sqlconn.SetMaxOpenConns(10) //设置连接池最大数
		sqlconn.SetMaxIdleConns(1)  //设置连接池中的最大闲置连接数。

		tracer := utils.GetTracer()
		//监听记录sql
		//db.Use(gormPlugin.New(tracer, gormPlugin.WithPeerAddr(host+":"+strconv.Itoa(port)), gormPlugin.WithSqlDBType(gormPlugin.MYSQL)))
		db.Use(gormPlugin.New(tracer,
			gormPlugin.WithPeerAddr(host+":"+strconv.Itoa(port)),
			gormPlugin.WithSqlDBType(gormPlugin.MYSQL),
			gormPlugin.WithQueryReport(), //打印sql语句
			gormPlugin.WithParamReport(), //打印sql的参数
		))

		return conn
	}

	return nil
}

//redis
func InitRedis(name string) *redis.Client {
	addr := conf.Get("hosts", name, "address").String("") + ":" + strconv.Itoa(conf.Get("hosts", name, "port").Int(0))
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

//func GetAuth(rdb *redis.Client, auth string) map[string]interface{} {
//	ctx := context.TODO()
//	ret := map[string]interface{}{}
//	data, err := rdb.Get(ctx, auth).Result()
//	if err == nil {
//		jsoniter.Unmarshal([]byte(data), &ret)
//		return ret
//	}
//
//	return nil
//}

//测试
func Test() {
	result := map[string]interface{}{}
	conn.Table("user").Find(&result)
	//fmt.Println(result)
}
