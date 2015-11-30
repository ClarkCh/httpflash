package httptask

import (
	"fmt"
	"strings"
	"time"
	"net/http"
	"net/url"
	"io/ioutil"
)

import (
	"httpflash/common"
)

func HttpTask() {
	fmt.Println("Ready")
	
	common.StartMutex.Lock()

	var client uint64

	var i uint64
	for i = 0; i < common.VCLIENT; i++ {
		go httpDo(common.VCLICK)
	}

	client = common.VCLIENT
	for _ = range common.StartChannel {
		client -= 1
		if client == 0 {
			fmt.Println("\nStart")
			close(common.StartChannel)
			common.StartMutex.Unlock()
		}
	}

	client = common.VCLIENT
	for _ = range common.StopChannel {
		client -= 1
		if client == 0 {
			close(common.StopChannel)
		}
	}

	fmt.Printf("==================================================================\n")
	fmt.Printf("All Client:              %16d\n", common.ResultPool.ClientNumber)
	fmt.Printf("All Click:               %16d\n", common.ResultPool.ClickNumber)
	fmt.Printf("All Seccess Click:       %16d\n", common.ResultPool.SeccessClickNumber)
	fmt.Printf("All False Click:         %16d\n", common.ResultPool.FalseClickNumber)

	if common.ResultPool.ClientNumber != 0 {
		t := common.ResultPool.SeccessTime/common.ResultPool.ClientNumber
		fmt.Printf("One Client Seccess Time: %16d(ns) %16d(ms)\n", t, t/1000000)
		t = common.ResultPool.FalseTime/common.ResultPool.ClientNumber
		fmt.Printf("One Client False Time:   %16d(ns) %16d(ms)\n", t, t/1000000)
	} else {
		fmt.Printf("One Client Seccess Time: %16d(ns) %16d(ms)\n", 0, 0)
		fmt.Printf("One Client False Time:   %16d(ns) %16d(ms)\n", 0, 0)
	}

	if common.ResultPool.SeccessClickNumber != 0 {
		t := common.ResultPool.SeccessTime/common.ResultPool.SeccessClickNumber
		fmt.Printf("One Click Seccess Time:  %16d(ns) %16d(ms)\n", t, t/1000000)
	} else {
		fmt.Printf("One Click Seccess Time:  %16d(ns) %16d(ms)\n", 0, 0)
	}

	if common.ResultPool.FalseClickNumber != 0 {
		t := common.ResultPool.FalseTime/common.ResultPool.FalseClickNumber
		fmt.Printf("One Click False Time:    %16d(ns) %16d(ms)\n", t, t/1000000)
	} else {
		fmt.Printf("One Click False Time:    %16d(ns) %16d(ms)\n", 0, 0)
	}

	fmt.Printf("False Connect:           %16d\n", common.ResultPool.FalseConnect)
}

func checkKeyValue(kvStr string, kvMap map[string]string) {
	if strings.Contains(kvStr, ",") {
		for _, str := range strings.Split(kvStr, ",") {
			buf := strings.Split(str, "=")
			key := buf[0]
			value := buf[1]

			kvMap[key] = value
		}
	} else {
		buf := strings.Split(kvStr, "=")
		key := buf[0]
		value := buf[1]

		kvMap[key] = value
	}
}

func httpDo(click uint64) {
	var (
		seccessClick int64
		falseClick   int64
		seccessTime  time.Duration
		falseTime    time.Duration
		falseConnect int64
	)

	client := &http.Client{}

	req, err := http.NewRequest(common.METHOD, common.URL, nil)
	if err != nil {
		panic(err)
	}

	if common.HEADER != "" {
		headerMap := make(map[string]string)
		checkKeyValue(common.HEADER, headerMap)
		for h := range headerMap {
			req.Header.Set(h, headerMap[h])
		}
	}

	if common.METHOD == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		if common.DATA != "" {
			dataMap := make(map[string]string)
			checkKeyValue(common.DATA, dataMap)
			data := url.Values{}
			for d := range dataMap {
				data.Add(d, dataMap[d])
			}

			req.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))
		}
	}

	if common.COOKIE != "" {
		req.Header.Set("Cookie", common.COOKIE)
	}

	
	common.StartChannel <- 1
	fmt.Print(".")
	common.StartMutex.RLock()
	common.StartMutex.RUnlock()

	var i uint64
	for i = 0; i < click; i++ {
		_req := req
		
		start := time.Now()
		resp, err := client.Do(_req)
		end := time.Now()
		if err != nil {
			falseConnect += 1
		} else {
			resp.Body.Close()
			statusCode := resp.StatusCode
			switch {
			case statusCode <= 399:
				seccessClick += 1
				seccessTime += end.Sub(start)
			case statusCode >= 400:
				falseClick += 1
				falseTime += end.Sub(start)
			}
		}
	}

	common.ResultPool.ClientNumber += 1
	common.ResultPool.ClickNumber += int64(common.VCLICK)

	common.ResultPool.SeccessClickNumber += seccessClick
	common.ResultPool.SeccessTime += seccessTime.Nanoseconds()

	common.ResultPool.FalseClickNumber += falseClick
	common.ResultPool.FalseTime += falseTime.Nanoseconds()

	common.ResultPool.FalseConnect += falseConnect

	common.StopChannel <- 1
}
