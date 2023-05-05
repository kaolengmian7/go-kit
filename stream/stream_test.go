package stream

import (
	"fmt"
	"log"
	"testing"
)

type Member struct {
	UserId  int64
	GroupId int64
	Wallet  int64
}

// 这是使用 Stream 包的一个 Demo
// 下面的例子中，我们将对 GroupId < 3 的星球的所有成员赠送 1 块钱。
func TestName(t *testing.T) {
	res, err := From(GetMembers).Filter(MemberWalletFilter).Map(WalletUpdate, WithWorkerNum(10)).Done(Report)
	if err != nil {
		log.Printf("%+v", err)
	}
	log.Println(res)
}

func GetMembers(dst chan<- interface{}) {
	members := []*Member{
		{1, 1, 0},
		{2, 1, 0},
		{3, 1, 0},
		{4, 1, 0},
		{5, 1, 0},
	}
	dst <- members
	members = []*Member{
		{1, 2, 0},
		{2, 2, 0},
		{3, 2, 0},
		{6, 2, 0},
		{7, 2, 0},
	}
	dst <- members
	members = []*Member{
		{1, 3, 2},
		{2, 3, 2},
		{8, 3, 2},
		{9, 3, 2},
		{10, 3, 2},
	}
	dst <- members
}

func MemberWalletFilter(val interface{}) bool {
	members, ok := val.([]*Member)
	if !ok {
		log.Printf("%+v", val)
	}
	for _, m := range members {
		if m.GroupId > 2 {
			return false
		}
	}
	return true
}

func WalletUpdate(val interface{}) interface{} {
	members, ok := val.([]*Member)
	if !ok {
		log.Printf("%+v", val)
	}
	for _, m := range members {
		m.Wallet++
	}
	return members
}

func Report(src <-chan interface{}) (result interface{}, err error) {
	total := 0
	for val := range src {
		members, ok := val.([]*Member)
		if !ok {
			log.Printf("%+v", val)
		}
		for _, m := range members {
			log.Printf("%+v", m)
		}
		total += len(members)
	}
	return fmt.Sprintf("数据处理完成，共计为 %d 位星友赠送 1 块钱", total), nil
}
