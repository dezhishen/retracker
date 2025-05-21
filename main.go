package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hekmon/transmissionrpc/v3"
)

func main() {
	// tm_host := os.Args[0]
	// tm_port := os.Args[1]
	// tm_username := os.Args[2]
	// tm_password := os.Args[3]
	// tm_url := fmt.Sprintf("http://%s:%s@%s:%s/transmission/rpc", tm_username, tm_password, tm_host, tm_port)
	fmt.Println("请输入你的Transmission的ip地址:")
	var tm_host string
	fmt.Scanln(&tm_host)
	fmt.Println("请输入你的Transmission的端口号:")
	var tm_port string
	fmt.Scanln(&tm_port)
	fmt.Println("请输入你的Transmission的用户名:")
	var tm_username string
	fmt.Scanln(&tm_username)
	fmt.Println("请输入你的Transmission的密码:")
	var tm_password string
	fmt.Scanln(&tm_password)
	var ssl bool
	fmt.Println("是否使用https连接？(y/n):")
	var ssl_input string
	fmt.Scanln(&ssl_input)
	if ssl_input == "y" || ssl_input == "Y" {
		ssl = true
	} else if ssl_input == "n" || ssl_input == "N" {
		ssl = false
	} else {
		fmt.Println("输入错误，默认使用http连接")
		ssl = false
	}
	var tm_url string
	if ssl {
		tm_url = fmt.Sprintf("https://%s:%s@%s:%s/transmission/rpc", tm_username, tm_password, tm_host, tm_port)
	} else {
		tm_url = fmt.Sprintf("http://%s:%s@%s:%s/transmission/rpc", tm_username, tm_password, tm_host, tm_port)
	}
	endpoint, err := url.Parse(tm_url)
	if err != nil {
		panic(err)
	}
	tmcli, err := transmissionrpc.New(endpoint, nil)
	if err != nil {
		panic(err)
	}
	//请输入你要替换的tracker原地址
	fmt.Println("请输入你要替换的tracker原地址:")
	var old_tracker string
	fmt.Scanln(&old_tracker)
	//请输入你要替换的tracker新地址
	fmt.Println("请输入你要替换的tracker新地址:")
	var new_tracker string
	fmt.Scanln(&new_tracker)
	// tbt获取所有的种子
	torrents, err := tmcli.TorrentGetAll(context.TODO())
	if err != nil {
		panic(err)
	}
	// 如果列表为空，打印并退出
	if len(torrents) == 0 {
		fmt.Println("没有种子")
		return
	}
	// 遍历所有的种子
	fmt.Println("正在查找需要替换的tracker...")
	for _, torrent := range torrents {
		// 获取种子的tracker
		trackers := torrent.Trackers
		// 遍历所有的tracker
		var new_trackers = make([]string, 0)
		need_replace := false
		for _, tracker := range trackers {
			// 如果tracker的url包含old_tracker，则替换为new_tracker
			if tracker.Announce == old_tracker {
				need_replace = true
				new_trackers = append(new_trackers, new_tracker)
			} else {
				new_trackers = append(new_trackers, tracker.Announce)
			}
		}
		if need_replace {
			fmt.Println("找到需要替换的tracker，正在替换...")
			err = tmcli.TorrentSet(context.TODO(), transmissionrpc.TorrentSetPayload{
				IDs:         []int64{*torrent.ID},
				TrackerList: new_trackers,
			})
			if err != nil {
				fmt.Println("替换失败:", err)
				continue
			}
			fmt.Println("替换成功")
		}
	}
	fmt.Println("全部替换完成")
}
