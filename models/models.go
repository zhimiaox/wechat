/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package models

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/zhimiao/wechat/common"
	"github.com/gomodule/redigo/redis"
	"github.com/influxdata/influxdb-client-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"time"
)

var (
	Redis    *zmRedis
	Mysql    *gorm.DB
	InfluxDB *zmInflux
)

type zmRedis struct {
	redis.Pool
}

type zmInflux struct {
	Client *influxdb.Client
}

type CommonMap map[string]interface{}

type ModelBase1 struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	UpdateTime time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	CreateTime time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}

// Start 初始化数据
func Start() {
	loadRedis()
	loadMysql()
	loadInfluxDB()
}

// 初始化influxDb
func loadInfluxDB() {
	InfluxDB = &zmInflux{}
}

// Setup Initialize the Redis instance
func loadRedis() {
	Redis = &zmRedis{
		redis.Pool{
			MaxIdle:     common.Config.Redis.MaxIdle,
			MaxActive:   common.Config.Redis.MaxActive,
			IdleTimeout: time.Duration(common.Config.Redis.IdleTimeOut) * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", common.Config.Redis.Host)
				if err != nil {
					return nil, err
				}
				if common.Config.Redis.Auth != "" {
					if _, err := c.Do("AUTH", common.Config.Redis.Auth); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

// Setup Initialize the Mysql instance
func loadMysql() {
	var err error
	Mysql, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		common.Config.Mysql.User,
		common.Config.Mysql.Password,
		common.Config.Mysql.Host,
		common.Config.Mysql.Database,
	))
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return common.Config.Mysql.TablePrefix + defaultTableName
	}
	Mysql.LogMode(true)
	Mysql.SingularTable(true)
	Mysql.DB().SetMaxIdleConns(10)
	Mysql.DB().SetMaxOpenConns(100)
	Mysql.Callback().Create().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {
		if !scope.HasError() {
			if createTimeField, ok := scope.FieldByName("CreateTime"); ok {
				if createTimeField.IsBlank {
					createTimeField.Set(time.Now())
				}
			}
			if modifyTimeField, ok := scope.FieldByName("UpdateTime"); ok {
				if modifyTimeField.IsBlank {
					modifyTimeField.Set(time.Now())
				}
			}
		}
	})
	Mysql.Callback().Update().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {
		if _, ok := scope.Get("gorm:update_column"); !ok {
			scope.SetColumn("UpdateTime", time.Now())
		}
	})
}

// Exists check a key
func (r *zmRedis) Exists(key string) bool {
	conn := r.Pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Delete delete a kye
func (r *zmRedis) Del(key string) (bool, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("DEL", key))
}

func (m *zmInflux) DB() *influxdb.Client {
	client, err := influxdb.New(common.Config.InfluxDB.Host, common.Config.InfluxDB.Token)
	if err != nil {
		logrus.Warn("InfluxDB初始化失败")
	}
	return client
}

func (m *zmInflux) Write(bucket string, metric ...influxdb.Metric) (err error) {
	conn := m.DB()
	defer conn.Close()
	_, err = conn.Write(context.Background(), bucket, common.Config.InfluxDB.Org, metric...)
	return
}

func (m *zmInflux) QueryToRaw(flux string) (raw []byte, err error) {
	conn := m.DB()
	defer conn.Close()
	data, err := conn.QueryCSV(context.Background(), flux, common.Config.InfluxDB.Org)
	if err != nil {
		return
	}
	raw, err = ioutil.ReadAll(data)
	if err != nil {
		return
	}
	return
}

func (m *zmInflux) QueryToArray(flux string) (result []map[string]interface{}, err error) {
	conn := m.DB()
	defer conn.Close()
	data, err := conn.QueryCSV(context.Background(), flux, common.Config.InfluxDB.Org)
	if err != nil {
		return
	}
	for data.Next() {
		rows := make(map[string]interface{})
		err = data.Unmarshal(rows)
		if data.Unmarshal(rows) == nil {
			result = append(result, rows)
		}
	}
	return
}

// HGet 获取一个值
func (r *zmRedis) HGet(cacheKey, key string) interface{} {
	conn := r.Pool.Get()
	defer conn.Close()
	var (
		data []byte
		err  error
	)
	if data, err = redis.Bytes(conn.Do("HGET", cacheKey, key)); err != nil {
		return nil
	}
	var reply interface{}
	if err = json.Unmarshal(data, &reply); err != nil {
		return nil
	}
	return reply
}

// HSet 设置一个值
func (r *zmRedis) HSet(cacheKey, key string, val interface{}) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	var data []byte
	if data, err = json.Marshal(val); err != nil {
		return
	}
	_, err = conn.Do("HSET", cacheKey, key, data)
	return
}

// Get 获取一个值
func (r *zmRedis) Get(key string) interface{} {
	conn := r.Pool.Get()
	defer conn.Close()
	var data []byte
	var err error
	if data, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		return nil
	}
	var reply interface{}
	if err = json.Unmarshal(data, &reply); err != nil {
		return nil
	}
	return reply
}

// Set 设置一个值
func (r *zmRedis) Set(key string, val interface{}, timeout time.Duration) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	var data []byte
	if data, err = json.Marshal(val); err != nil {
		return
	}
	if timeout == 0 {
		_, err = conn.Do("SET", key, data)
		return
	}
	_, err = conn.Do("SETEX", key, int64(timeout/time.Second), data)
	return
}

// IsExist 判断key是否存在
func (r *zmRedis) IsExist(key string) bool {
	conn := r.Pool.Get()
	defer conn.Close()

	a, _ := conn.Do("EXISTS", key)
	i := a.(int64)
	if i > 0 {
		return true
	}
	return false
}

// Delete 删除
func (r *zmRedis) Delete(key string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", key); err != nil {
		return err
	}

	return nil
}
