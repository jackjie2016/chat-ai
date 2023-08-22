package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	url := "https://gw.newrank.cn/api/xd/xdnphb/nr/cloud/douyin/user/trend/performance"
	method := "POST"

	payload := strings.NewReader(`{"uid":"MS4wLjABAAAACirq2YgCDHFt4JJLwl1l5Hj4WpThkGSm8uKQJY7a2hU","dateType":"2","hasChart":false,"startTime":"","endTime":""}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("gw-c-v", "10000")
	req.Header.Add("n-token", "9116298d52d64bbfb2bafa92267f74f2")
	req.Header.Add("Cookie", "__bid_n=187b6c5e0426864765ef3f; Hm_lvt_a19fd7224d30e3c8a6558dcb38c4beed=1686821239,1686894846,1687659997,1688009711; token=BB9E18BC95F94BCD9596A66DA1083FA9; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%22nr_17npvqfsi%22%2C%22first_id%22%3A%221877af1387c445-0769c56cee830e-9196d2c-2073600-1877af1387d47f%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%2C%22%24latest_referrer%22%3A%22%22%2C%22%24latest_utm_source%22%3A%22websites%22%2C%22%24latest_utm_medium%22%3A%22xd_banner%22%2C%22%24latest_utm_campaign%22%3A%22%E6%96%B0%E6%A6%9C%E6%9C%89%E6%95%B0%22%2C%22%24latest_utm_content%22%3A%225.12%20%E8%AF%9D%E9%A2%98%E7%83%AD%E6%A6%9C%E4%B8%8A%E7%BA%BF%22%2C%22%24latest_utm_term%22%3A%22%E6%96%B0%E6%8A%96-%E7%9F%AD%E8%A7%86%E9%A2%91-%E7%83%AD%E9%97%A8%E8%AF%9D%E9%A2%98%22%7D%2C%22%24device_id%22%3A%221877af1387c445-0769c56cee830e-9196d2c-2073600-1877af1387d47f%22%7D; acw_tc=707c9fc416885197883042114e4658ec617da7bcb1870146adcb5fdd0ceb0f; NR_MAIN_SOURCE_RECORD={\"locationSearch\":\"?ruuid")
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "gw.newrank.cn")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
