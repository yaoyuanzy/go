package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/oschwald/geoip2-golang"
	"github.com/otiai10/copy"
)

type IntMapOrderedItem struct {
	Key   string
	Value int
}

type IntMapOrdered []*IntMapOrderedItem

func NewIntMapOrdered(m map[string]int) IntMapOrdered {
	im := make(IntMapOrdered, len(m))
	index := 0
	for key, value := range m {
		im[index] = &IntMapOrderedItem{key, value}
		index++
	}
	return im
}

func (m IntMapOrdered) Get(key string) (int, error) {
	for _, item := range m {
		if item.Key == key {
			return item.Value, nil
		}
	}
	return 0, fmt.Errorf("key not found")
}

func (m IntMapOrdered) Set(key string, value int) {
	for _, item := range m {
		if item.Key == key {
			item.Value = value
			return
		}
	}
	m = append(m, &IntMapOrderedItem{key, value})
}

func (m IntMapOrdered) Delete(key string) {
	index := -1
	for idx, item := range m {
		if item.Key == key {
			index = idx
			break
		}
	}
	if index != -1 {
		m = append(m[:index], m[index+1:]...)
	}
}

func (m IntMapOrdered) Len() int {
	return len(m)
}

func (m IntMapOrdered) Less(i, j int) bool {
	return m[i].Value < m[j].Value
}

func (m IntMapOrdered) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func SortMap(m IntMapOrdered, order string) IntMapOrdered {
	if order == "desc" {
		sort.Sort(sort.Reverse(m))
	} else {
		sort.Sort(m)
	}
	return m
}

func Json(v interface{}) string {
	if bytes, err := json.Marshal(v); err == nil {
		return string(bytes)
	}
	return "{}"
}

func main() {
	// 定义命令行参数, 接收web访问日志文件路径
	path := flag.String("path", "", "access log")
	dir := flag.String("dir", "", "report log")

	// 定义命令行使用说明
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", "loganalysis")
		flag.PrintDefaults()
	}

	// 解析命令行参数
	flag.Parse()

	// 检查结果目录
	if *dir == "" {
		*dir = fmt.Sprintf("reports/report_%d", time.Now().Unix())
	}
	if _, err := os.Stat(*dir); err == nil {
		fmt.Printf("结果目录`%s`已存在\n", *dir)
		os.Exit(-1)
	}

	// 根据参数获取路径信息
	file, err := os.Open(*path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("路径`%s`不存在\n", *path)
		} else {
			fmt.Println(err)
		}
		os.Exit(-1)
	} else {
		defer file.Close()
		if info, _ := file.Stat(); info.IsDir() {
			fmt.Printf("路径`%s`不能为目录\n", *path)
		} else {
			// 定义每天统计信息
			var (
				days         map[string]bool //定义日期
				hitTotal     int             //总的点击数量
				bytesTotal   uint64          //总的字节数量
				vistorsTotal map[string]int  //每个IP的访问次数
				statusTotal  map[string]int  // 每个状态码出现次数
			)
			days = make(map[string]bool)
			vistorsTotal = make(map[string]int)
			statusTotal = make(map[string]int)

			hitDays := make(map[string]int)                // 每天访问量
			bytesDays := make(map[string]uint64)           //每天流量大小
			vistorsDays := make(map[string]map[string]int) //每天每个IP访问次数
			statusDays := make(map[string]map[string]int)  //每天每个状态码出现次数

			// 使用带缓冲区的扫描器，每次读取一行内容
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				nodes := strings.Split(line, " ")
				if len(nodes) < 12 {
					continue
				}
				// ip datetime method url status bytes
				// fmt.Printf("%q, %q, %q, %q, %q, %q\n", nodes[0], nodes[3][1:], nodes[5][1:], nodes[6], nodes[8], nodes[9])

				// 时间处理
				logTime, _ := time.Parse("02/Jan/2006:15:04:05", nodes[3][1:])
				logDay := logTime.Format("2006-01-02")

				// 设置日期
				days[logDay] = true

				// 点击数量统计
				hitTotal++
				hitDays[logDay]++

				// 使用流量统计
				if b, err := strconv.ParseUint(nodes[9], 10, 64); err == nil {
					bytesTotal += b
					bytesDays[logDay] += b
				}

				// 初始化每天访问量和状态码统计map
				if _, exists := vistorsDays[logDay]; !exists {
					vistorsDays[logDay] = make(map[string]int)
				}
				if _, exists := statusDays[logDay]; !exists {
					statusDays[logDay] = make(map[string]int)
				}

				// 访问IP统计
				vistorsTotal[nodes[0]]++
				vistorsDays[logDay][nodes[0]]++

				// 状态码统计
				statusTotal[nodes[8]]++
				statusDays[logDay][nodes[8]]++
			}
			if err := scanner.Err(); err != nil {
				fmt.Println(err)
			} else {
				err := copy.Copy("tpl", *dir) //复制模板
				if err != nil {
					fmt.Println(err)
					os.Exit(-2)
				}
				// 定义过滤器函数
				funcMaps := template.FuncMap{
					"filesizeformat": humanize.Bytes,
					"sortmap":        SortMap,
					"json":           Json,
				}

				// 定义模板
				tmpl := template.New("index.html").Funcs(funcMaps)
				tmpl = template.Must(tmpl.ParseFiles(filepath.Join("tpl", "index.html")))

				// 定义模板输出文件
				file, err := os.Create(filepath.Join(*dir, "index.html"))
				if err != nil {
					fmt.Println(err)
					os.Exit(-2)
				}

				defer file.Close()

				// 带缓冲输出
				writer := bufio.NewWriter(file)
				defer writer.Flush()

				dayKeys := make([]string, len(days))
				index := 0
				for key := range days {
					dayKeys[index] = key
					index++
				}

				// 打开地址库
				geoip, err := geoip2.Open("db/GeoLite2-City.mmdb")
				if err != nil {
					fmt.Println(err)
					os.Exit(-2)
				}
				defer geoip.Close() //延迟关闭

				regionTotal := make(map[string]int)
				regionLocation := make(map[string][2]float64)
				for vistor, cnt := range vistorsTotal { //遍历访问IP地址
					ip := net.ParseIP(vistor) //解析为net.IP类型
					if ip == nil {
						fmt.Println(vistor)
						continue
					}
					record, err := geoip.City(ip) //查询IP信息
					if err != nil {
						fmt.Println(vistor, err)
						continue
					}

					//只显示国内IP地址
					if record.Country.Names["zh-CN"] != "中国" {
						continue
					}

					name := fmt.Sprintf("%s/%s", record.Country.Names["zh-CN"], record.City.Names["zh-CN"])
					regionTotal[name] += cnt
					if _, exists := regionLocation[name]; !exists {
						regionLocation[name] = [2]float64{record.Location.Longitude, record.Location.Latitude}
					}
				}

				err = tmpl.Execute(writer, struct {
					Days         []string
					HitTotal     int
					BytesTotal   uint64
					VistorsTotal IntMapOrdered
					StatusTotal  map[string]int

					HitDays     map[string]int
					BytesDays   map[string]uint64
					VistorsDays map[string]map[string]int
					StatusDays  map[string]map[string]int

					RegionTotal    map[string]int
					RegionLocation map[string][2]float64
				}{
					Days:           dayKeys,
					HitTotal:       hitTotal,
					BytesTotal:     bytesTotal,
					VistorsTotal:   NewIntMapOrdered(vistorsTotal),
					StatusTotal:    statusTotal,
					HitDays:        hitDays,
					BytesDays:      bytesDays,
					VistorsDays:    vistorsDays,
					StatusDays:     statusDays,
					RegionTotal:    regionTotal,
					RegionLocation: regionLocation,
				})

				if err != nil {
					fmt.Println(err)
					os.Exit(-2)
				}
				fmt.Println("Success")
			}
		}
	}

}
