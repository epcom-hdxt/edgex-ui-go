/*******************************************************************************
 * Copyright © 2020-2021 VMware, Inc. All Rights Reserved.
 * edgex1.3.0版本未能实现批量删除redis key功能
 * 因此实现一个方法，用于定时任务调用，搜索出redis中keys event:readings* set
 * 执行批量删除
 *
 *******************************************************************************/

package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
)

/**
 * event 需要清除的key
 * 1.event:created
 * 2.event:device*  set 循环清除
 * 3.event:pushed
 * 4.event:readings* set 循环清除 并 清除set里面的key
 * 5.event
 */
func ScrubEventReadingAndDel(w http.ResponseWriter, r *http.Request) {

	//TODO 等孩子睡了再书写
	fmt.Print(1)
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Printf(err.Error())
	}
	defer conn.Close()
	// fmt.Print("conn ok!")
	// conn.Send("MULTI")

	valset, err := redis.Strings(conn.Do("KEYS", "event:readings*"))
	if err != nil {
		log.Printf(err.Error())
	}
	// fmt.Println(len(valset), err)

	for i, _ := range valset {

		// 读取指定zset
		set_map, err := redis.StringMap(conn.Do("ZRANGE", valset[i], 0, -1, "withscores"))
		if err != nil {
			log.Println("redis get failed:", err)
		} else {
			// fmt.Println("Get " + valset[i] + " 长度:" + strconv.Itoa(len(set_map)))
		}

		for setmap := range set_map {
			//fmt.Printf("set: %v %v\n", setmap, set_map[setmap])
			// fmt.Println(conn.Do("DEL", setmap))
			conn.Do("DEL", setmap) //删除set内部 key
		}
		conn.Do("DEL", valset[i]) //删除set本身
		// fmt.Println("删除" + valset[i] + "ok ----" + "第" + strconv.Itoa(i) + "个")

	}
	conn.Do("DEL", "event:created") //删除event:created
	conn.Do("DEL", "event:pushed")  //删除event:created
	conn.Do("DEL", "event")         //删除event:created

	//循环删除event:device*
	event_device, err := redis.Strings(conn.Do("KEYS", "event:device*"))
	if err != nil {
		log.Println("redis get failed:", err)
	} else {
		for i, _ := range event_device {
			// fmt.Println(event_device[i])
			conn.Do("DEL", event_device[i])
		}
	}
	log.Println("redis scrub event and delete success!")
	w.Write([]byte(nil))
}
