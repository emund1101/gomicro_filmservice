package handler

import (
	"context"
	pb "films/service_film/proto"
	"films/utils"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"go-micro.dev/v4/client"
	"gorm.io/gorm"
	"time"
)

type FilmService struct {
	Db        *gorm.DB
	Rdb       *redis.Client
	GRPClient client.Client
}

type Filmlist struct {
	Name     string
	Date     string
	End_date string
	Desc     string
	Poster   string
	Video    string
	Score    int
}

type Filmhallinfo struct {
	Name  string
	Price float32
	Seats []map[string]interface{}
}

func (e *FilmService) GetAllList(ctx context.Context, req *pb.FilmListRequest, rsp *pb.Respone) error {
	//调用验证登录的微服务
	if rs := e.getAuth(req.Auth); rs == false {
		byte, _ := jsoniter.Marshal(utils.Json_return{Code: 0, Msg: "授权码失效", Data: nil}) //输出数据
		rsp.Respone = byte
		return nil
	}
	var rst []Filmlist
	e.Db.Table("films").Where("end_date>=?", time.Now().Format("2006-01-02")).Order("end_date desc").Find(&rst)
	//	fmt.Println(rst)
	byte, _ := jsoniter.Marshal(utils.Json_return{Code: 1, Msg: "成功", Data: rst})
	rsp.Respone = byte
	return nil
}

//获取参数下该电影等的数据场次信息
func (e *FilmService) GetInfoList(ctx context.Context, req *pb.FilmListRequest, rsp *pb.Respone) error {
	//调用验证登录的微服务
	if rs := e.getAuth(req.Auth); rs == false {
		byte, _ := jsoniter.Marshal(utils.Json_return{Code: 0, Msg: "授权码失效", Data: nil}) //输出数据
		rsp.Respone = byte
		return nil
	}

	var rst map[string]string
	//select f.name,group_concat(fi.id) as sid,c.name as cname,c.address,c.mobile,group_concat(fi.time) as vtime,group_concat(fi.price) as price from films as f join films_info as fi on f.id=fi.fid join cinema as c on c.id=fi.cid
	//where DATE_FORMAT(fi.time, '%Y-%m-%d') = '2022-05-16' and fi.status=1 group by c.id;
	field := " f.name,group_concat(fi.id) as sid,c.name as cname,c.address,c.mobile,group_concat(fi.time) as vtime,group_concat(fi.price) as price"

	//读取电影数据
	gorm := e.Db.Table("films as f").Joins("join films_info fi on f.id=fi.fid").Joins("join cinema as c on c.id=fi.cid").Select(field).Where("c.status=?", 1).Where("fi.status=?", 1)
	if req.Name != "" {
		gorm.Where("f.name like %?%", req.Name)
	}
	if req.Cinema != "" {
		gorm.Where("c.name like %?%", req.Cinema)
	}
	if req.Time != "" {
		gorm.Where("DATE_FORMAT(fi.time,'%Y-%m-%d') = ?", req.Time)
		//	gorm.Where("f.date >=? and f.end_date<=?", req.Time, req.Time)
	}

	gorm.Find(&rst)

	byte, err := jsoniter.Marshal(utils.Json_return{Code: 1, Msg: "成功", Data: rst}) //输出数据
	if err == nil {
		rsp.Respone = byte
	}

	return nil
}

//获取位置信息
func (e *FilmService) GetHallInfo(ctx context.Context, req *pb.FilmInfoRequest, rsp *pb.Respone) error {
	//调用验证登录的微服务
	if rs := e.getAuth(req.Auth); rs == false {
		byte, _ := jsoniter.Marshal(utils.Json_return{Code: 0, Msg: "授权码失效", Data: nil}) //输出数据
		rsp.Respone = byte
		return nil
	}

	rs := []map[string]interface{}{}
	field := "s.id,s.line,s.no,s.status,fi.price" //字段名
	e.Db.Table("films_info as fi").Joins("join seat as s on s.hall=fi.hallid").Where("fi.id = ? ", req.FilmSid).Select(field).Find(&rs)

	list, err := e.Rdb.HGetAll(ctx, string(req.FilmSid)).Result()
	if err == nil {
		for _, item := range rs {
			for k, l := range list {
				if k == item["id"] {
					item["status"] = l //redis 的值，2,未付款锁定，2 是已付款锁定
				}
			}
		}
	}

	data := Filmhallinfo{}
	data.Seats = rs

	byte, err := jsoniter.Marshal(utils.Json_return{Code: 1, Msg: "成功", Data: data}) //输出数据
	rsp.Respone = byte
	return nil
}

//验证登录授权码
func (e *FilmService) getAuth(auth string) bool {
	var rsp map[string][]byte
	request := map[string]string{"auth": auth}
	greq := e.GRPClient.NewRequest("user", "UserService.Authorize", request, client.WithContentType("application/json"))

	ctx := utils.GetContext()
	//注入链路context
	//	plugin_micro.NewCall(utils.GetTracer(), ctx, greq)
	//	nctx := utils.GetContext()
	//	if err := e.GRPClient.Call(nctx, greq, &rsp); err == nil {
	
	if err := e.GRPClient.Call(ctx, greq, &rsp); err == nil {
		rst := utils.Json_return{}
		jsoniter.Unmarshal(rsp["respone"], &rst)
		if rst.Code == 1 {
			return true
		}

	}
	
// 	e.GRPClient.Call(ctx, greq, &rsp)
// 	rst := utils.Json_return{}
// 	jsoniter.Unmarshal(rsp["respone"], &rst)
// 	if rst.Code == 1 {
// 		return true
// 	}

	return false
}
