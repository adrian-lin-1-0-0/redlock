package main

import (
	"log"
	"time"

	"github.com/adrian-lin-1-0-0/redlock"
)

func main() {

	rl := redlock.New(
		&redlock.Options{
			Addr: "localhost:6379",
		},
		&redlock.Options{
			Addr: "localhost:6380",
		},
	)

	m := rl.NewMutex("test1", 10*time.Second)
	err := m.Lock()
	if err != nil {
		log.Default().Println(err)
		m.Unlock()
		return
	}

	mysql := NewMysql()
	data := mysql.GetTestById(1)
	version := data.Version
	newVersion := version + 1
	mysql.UpdateTest.
		Where("id", 1).
		Where("version", version).
		Set("version", newVersion).
		Set("value", "new value").
		Exec()

	m.Unlock()
}

func NewMysql() *Mysql {
	return &Mysql{
		UpdateTest: UpdateTest{
			Id:      1,
			Version: 1,
			Value:   "test",
		},
	}
}

type Mysql struct {
	UpdateTest UpdateTest
}

type UpdateTest struct {
	Id      int
	Version int
	Value   string
}

func (u *UpdateTest) Where(key string, value interface{}) *UpdateTest {
	return u
}

func (u *UpdateTest) Set(key string, value interface{}) *UpdateTest {
	if key == "version" {
		u.Version = value.(int)
	}

	if key == "value" {
		u.Value = value.(string)
	}

	return u
}

func (m Mysql) GetTestById(id int) *UpdateTest {
	return &m.UpdateTest
}

func (u *UpdateTest) Exec() {

}
