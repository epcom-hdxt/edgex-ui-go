/*******************************************************************************
 * Copyright © 2020-2021 VMware, Inc. All Rights Reserved.
 * edgex1.3.0版本未能实现批量删除redis key功能
 * 因此实现一个方法，用于定时任务调用，搜索出redis中keys
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
 * 1.reading:created
 * 2.reading:device*  set 循环清除
 * 3.reading:name* set 循环清除 并 清除set里面的key
 * 4.reading
 */
func ScrubReadingAndDel(w http.ResponseWriter, r *http.Request) {

	//TODO 等孩子睡了再书写
	fmt.Print(1)
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Print("conn ok!")

	conn.Do("DEL", "reading:created") //删除reading:created
	conn.Do("DEL", "reading")         //删除reading

	//循环删除reading:device*
	reading_device, err := redis.Strings(conn.Do("KEYS", "reading:device*"))
	if err != nil {
		log.Printf(err.Error())
	} else {
		for i, _ := range reading_device {
			// fmt.Println(reading_device[i])
			conn.Do("DEL", reading_device[i])
		}
	}

	//循环删除reading:name*
	reading_name, err := redis.Strings(conn.Do("KEYS", "reading:name*"))
	if err != nil {
		// fmt.Println("redis get event_device failed:", err)
		log.Printf(err.Error())
	} else {
		for i, _ := range reading_name {
			// fmt.Println(reading_name[i])
			conn.Do("DEL", reading_name[i])
		}
	}
	log.Println("redis scrub readingss and delete success!")
	w.Write([]byte(nil))
}
