package redis

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestDemo1(t *testing.T) {
	path, err := filepath.Abs("./conf.yaml")
	if err != nil {
		panic(err)
	}
	fileRead, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	conf := &Config{}
	if err = yaml.Unmarshal(fileRead, &conf); err != nil {
		panic(err)
	}

	ctx := context.Background()
	client := New(conf)
	defer client.Close()

	//_ = client.Set(ctx, conf.Name, "aa", time.Second*10)

	//result, err := client.Get(ctx, conf.Name).Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(result)

	conn := client.Conn(ctx)
	//defer conn.Close()
	result1, err := conn.Get(ctx, conf.Name).Result()
	if err != nil {
		panic(err)
	}
	//bytes, _ := json.Marshal(client.PoolStats())
	//fmt.Println(string(bytes), result1)
	fmt.Println(result1)
}
