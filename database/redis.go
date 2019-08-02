package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	opt, err := redis.ParseURL("redis://:@127.0.0.1:7000/0")
	if err != nil {
		panic(err)
	}
	
	if client == nil {
		client = redis.NewClient(opt)
	}
	ping()
}

// DbPerson person
type DbPerson struct {
	Age  int
	Sex  int
	Desc string
	name string
}

func (p *DbPerson) writeTo(client redis.Cmdable) {
	// hash set
	j, err := json.Marshal(p)
	if err != nil {
		fmt.Println("fail:", err)
	}

	err = client.HSet("person", p.name, string(j)).Err()
	if err != nil {
		fmt.Println("fail:", err)
	}
}

func (p *DbPerson) decodeFromStrCmd(r *redis.StringCmd) {
	if err := r.Err(); err != nil {
		fmt.Println("fail:", err)
	} else {
		if b, err := r.Bytes(); err != nil {
			fmt.Println("fail:", err)
		} else {
			json.Unmarshal(b, p)
		}
	}
}

func (p *DbPerson) readFrom(client redis.Cmdable) {
	p.decodeFromStrCmd(client.HGet("person", p.name))
}

func (p DbPerson) String() string {
	return fmt.Sprintf("name: %s, age: %d, sex: %d, desc: %s", p.name, p.Age, p.Sex, p.Desc)
}

// RedisExample example
func RedisExample() {
	p := DbPerson{30, 1, "一个普通人", ""}

	// withPipelining(func(db redis.Cmdable) {
	// 	for i := 0; i < 100; i++ {
	// 		p.Age = i
	// 		p.name = fmt.Sprintf("gaoyang_%d", i)
	// 		p.writeTo(db)
	// 	}
	// })

	// withoutPipelining(func(db redis.Cmdable) {
	// 	for i := 0; i < 100; i++ {
	// 		p.Age = i
	// 		p.name = fmt.Sprintf("gaoyang_%d", i)
	// 		p.writeTo(db)
	// 	}
	// })

	results := []*redis.StringCmd{}
	withPipelining(func(db redis.Cmdable) {
		for i := 0; i < 100; i++ {
			p.Age = i
			p.name = fmt.Sprintf("gaoyang_%d", i)
			results = append(results, client.HGet("person", p.name))
		}
	}, func(cmds []redis.Cmder, err error) {
		if err != nil {
			fmt.Println(err)
			return
		}

		p := DbPerson{}
		for i := range results {
			args := results[i].Args()
			p.name = args[len(args)-1].(string)
			p.decodeFromStrCmd(results[i])
			fmt.Println(p)
		}
	})
}

func withoutPipelining(content func(client redis.Cmdable)) {
	now := time.Now()
	content(client)
	fmt.Println("all done. ", time.Now().Sub(now).Seconds(), "sec")
}

func withPipelining(content func(client redis.Cmdable), result func([]redis.Cmder, error)) {
	pip := client.Pipeline()
	now := time.Now()

	content(pip)
	if result != nil {
		result(pip.Exec())
	}

	fmt.Println("pipe done. ", time.Now().Sub(now).Seconds(), "sec")
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

func clusterExample() {

}