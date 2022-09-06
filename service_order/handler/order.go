package handler

import (
	"context"
	pb "films/service_order/proto"
	"films/utils/orm"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

type OrderService struct {
	Db    *gorm.DB
	Redis *redis.Client
}

type order struct {
	Id          int
	Orderno     string
	Status      int
	User_id     int64
	Create_time string
	Num         int
	Total       float32
}

type orderInfo struct {
	Id      int
	Oid     int
	Film_id int64
	Seat_id int
	Qrcode  string
	Price   float32
}

//订单状态
const (
	OrderStatus_nopay  = 0 //未支付
	OrderStatus_pay    = 1 //已支付
	OrderStatus_cancel = 2 //取消
)

//下单预订
func (e *OrderService) Buy(ctx context.Context, req *pb.BuyRequest, rsq *pb.Respone) error {
	var err error //记录错误
	//判断数据
	seats := strings.Split(req.Seats, ",")

	//锁定座位,只保留10分钟,10分钟后删除
	for _, item := range seats {
		cmd := e.Redis.HSetNX(ctx, string(req.FilmSid), item, 1) //filmid,seatid,
		if err := cmd.Err(); err != nil {
			return err
		}
	}

	rand.Seed(time.Now().UnixNano())
	data := order{
		Orderno:     time.Now().Format("20060102150405") + string(1000+rand.Intn(8999)),
		Status:      OrderStatus_nopay,
		User_id:     req.Userid,
		Create_time: time.Now().Format("2006-01-02 15:04:05"),
		Num:         len(seats),
	}
	e.Db.Table("order").Create(&data)
	datainfo := make([]orderInfo, 0, cap(seats))
	for seat := range seats {
		datainfo = append(datainfo, orderInfo{
			Oid:     data.Id,
			Film_id: req.FilmSid,
			Seat_id: int(seat),
		})
	}
	e.Db.Table("order_info").Create(&datainfo)

	var rst orm.Json_return
	if err == nil {
		rst = orm.Json_return{1, "操作成功", nil}
	} else {
		rst = orm.Json_return{0, err.Error(), nil}
	}

	json, _ := jsoniter.Marshal(rst)
	rsq.Respone = json

	return nil
}

//付款
func (e *OrderService) Pay(ctx context.Context, req *pb.PayRequest, rsq *pb.Respone) error {
	return nil
}

func (e *OrderService) Refund(ctx context.Context, req *pb.PayRequest, rsq *pb.Respone) error {
	return nil
}

func (e *OrderService) GetInfo(ctx context.Context, req *pb.InfoRequest, rsq *pb.Respone) error {
	return nil
}
