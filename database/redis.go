package database

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

var client *redis.Client

// DbPerson person
type DbPerson struct {
	Age  int
	Sex  int
	Desc string
}

func init() {
	opt, err := redis.ParseURL("redis://:shine@192.168.1.4:6379/0")
	if err != nil {
		panic(err)
	}

	if client == nil {
		client = redis.NewClient(opt)
	}
}

func redisExample() {
	p := DbPerson{30, 1, "一个普通人"}
	for i := 0; i < 100000; i++ {
		p.Age = i
		setPerson(fmt.Sprintf("gaoyang_%d", i), &p)
	}
}

func ping() {
	// ping
	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println("fail:", err)
	} else {
		fmt.Println("ping:", pong)
	}
}

func setPerson(name string, p *DbPerson) {
	// hash set
	j, err := json.Marshal(p)
	if err != nil {
		fmt.Println("fail:", err)
	}

	err = client.HSet("person", name, string(j)).Err()
	if err != nil {
		fmt.Println("fail:", err)
	}
}

func getPerson(name string) *DbPerson {
	r := client.HGet("person", name)
	if err := r.Err(); err != nil {
		fmt.Println("fail:", err)
	} else {
		if b, err := r.Bytes(); err != nil {
			fmt.Println("fail:", err)
		} else {
			p := DbPerson{}
			json.Unmarshal(b, &p)
			return &p
		}
	}
	return nil
}
