package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

type job struct {
	idx int
	len int
}

type Jobs []job

func (jb Jobs) Len() int {
	return len(jb)
}

func (jb Jobs) Swap(i, j int) {
	jb[i], jb[j] = jb[j], jb[i]
}

func (jb Jobs) Less(i, j int) bool {
	return jb[i].len < jb[j].len
}

var seed = flag.Int64("s", 0, "随机种子")
var jobs = flag.Int("j", 3, "作业数量")
var jlist = flag.String("l", "", "提供一个以逗号分隔的运行时间列表，而不是随机长度作业")
var maxlen = flag.Int("m", 10, "作业的最大长度")
var policy = flag.String("p", "FIFO", "调度策略: SJF, FIFO, RR")
var quantum = flag.Int("q", 1, "RR策略的时间片长度")
var solve = flag.Bool("c", false, "计算答案")

func main() {
	flag.Parse()

	// 输出提示
	fmt.Printf("调度策略为： %s\n", *policy)
	if *jlist == "" {
		fmt.Printf("调度任务数为： %d\n", *jobs)
		fmt.Printf("调度任务作业最大长度为： %d\n", *maxlen)
		fmt.Printf("调度种子为： %d\n", *seed)
	} else {
		fmt.Printf("调度任务为： %s\n", *jlist)
	}

	// 解析任务
	fmt.Printf("以下是作业列表，其中包含每个作业的运行时间\n")
	var joblist Jobs
	if *jlist != "" {
		jobl := strings.Split(*jlist, ",")
		joblist = make([]job, len(jobl))
		for idx, item := range jobl {
			joblen, _ := strconv.Atoi(item)
			joblist[idx] = job{
				idx: idx,
				len: joblen,
			}
		}
	} else {
		rander := rand.New(rand.NewSource(*seed))
		joblist = make([]job, *jobs)
		for i := 0; i < *jobs; i++ {
			joblist[i] = job{
				idx: i,
				len: rander.Intn(*maxlen) + 1,
			}
		}
	}
	for _, v := range joblist {
		fmt.Printf("\tJob %d %d\n", v.idx, v.len)
	}

	if !*solve {
		fmt.Println("计算每个作业的周转时间、响应时间和等待时间。")
		fmt.Println("完成后，使用相同的参数再次运行此程序。")
		fmt.Println("但是使用-c，这将为您提供答案。您可以使用-s<somenumber>或您自己的工作列表（例如-l 10、15、20），改变不同的题型")
		return
	}

	fmt.Println("** Solutions **")

	switch *policy {
	case "SJF":
		sjf(joblist)
		return
	case "FIFO":
		fifo(joblist)
		return
	case "RR":
		rr(joblist, *quantum)
		return
	default:
		fmt.Println("策略有误")
		return
	}
}

// 最短任务优先
func sjf(joblist Jobs) {
	sort.Sort(joblist)
	fifo(joblist)
}

// fifo 先入先出
func fifo(joblist Jobs) {
	thetime := 0
	fmt.Println("调用栈：")
	for _, job := range joblist {
		fmt.Printf("[ time %3d ] Run job %d for %3d secs ( DONE at %3d )", thetime, job.idx, job.len, thetime+job.len)
		thetime += job.len
	}

	fmt.Println("\n统计")
	t := 0
	turnaroundSum := 0
	waitSum := 0
	responseSum := 0
	for _, job := range joblist {
		response := t
		turnaround := t + job.len
		wait := t
		fmt.Printf("\tJob %3d -- 响应时间: %3d  周转时间 %3d  等待时间 %3d\n", job.idx, response, turnaround, wait)
		responseSum += response
		turnaroundSum += turnaround
		waitSum += wait
		t += job.len
	}
	fmt.Printf("\n平均 -- 响应时间: %3d  周转时间 %3d  等待时间 %3d\n", responseSum/len(joblist), turnaroundSum/len(joblist), waitSum/len(joblist))
}

// 轮转任务
func rr(joblist Jobs, quantum int) {
	thetime := 0
	turnaround := make(map[int]int, len(joblist))
	response := make(map[int]int, len(joblist))
	lastrun := make(map[int]int, len(joblist)) // 上次运行时间
	wait := make(map[int]int, len(joblist))    // 等待时间

	jobcount := len(joblist)
	runlist := make([]job, 0, 100)
	for _, item := range joblist {
		lastrun[item.idx] = 0
		response[item.idx] = 0
		turnaround[item.idx] = 0
		wait[item.idx] = 0
		runlist = append(runlist, item)
	}

	fmt.Println("调用栈：")
	for jobcount > 0 {
		jobitem := runlist[0]
		runlist = runlist[1:]
		runtime := jobitem.len
		currwait := thetime - lastrun[jobitem.idx]
		wait[jobitem.idx] += currwait

		var runfor int
		if runtime > quantum {
			runtime -= quantum
			runfor = quantum
			runlist = append(runlist, job{
				idx: jobitem.idx,
				len: runtime,
			})
			fmt.Printf("\t[ time %3d ] Run job %3d for %3d secs\n", thetime, jobitem.idx, runfor)
		} else {
			runfor = runtime
			turnaround[jobitem.idx] = thetime + runfor
			jobcount--
		}
		thetime += runfor
		lastrun[jobitem.idx] = thetime
	}

	fmt.Println("\n统计")
	turnaroundSum := 0
	waitSum := 0
	responseSum := 0
	for _, item := range joblist {
		turnaroundSum += turnaround[item.idx]
		responseSum += response[item.idx]
		waitSum += wait[item.idx]
		fmt.Printf("\tJob %3d -- 响应时间: %3d  周转时间 %3d  等待时间 %3d\n", item.idx, response[item.idx], turnaround[item.idx], wait[item.idx])
	}
	fmt.Printf("\n平均 -- 响应时间: %3d  周转时间 %3d  等待时间 %3d\n", responseSum/len(joblist), turnaroundSum/len(joblist), waitSum/len(joblist))
}
