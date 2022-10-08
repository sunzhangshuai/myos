package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// Job2 任务
type Job2 struct {
	num     int // 任务编号
	runTime int // 执行时间
	tickets int // 彩票数
}

var seed2 = flag.Int64("s", 0, "指定随机种子")
var jobnum = flag.Int("j", 3, "系统中的作业数量")
var jlist2 = flag.String("l", "", "不使用随机作业，而是提供以逗号分隔的运行时间和彩票数列表。(例如10:100,20:100将有两个作业，其运行时间分别为10和20，每个作业具有100张彩票)")
var maxlen2 = flag.Int("m", 10, "工作的最大长度")
var maxticket = flag.Int("T", 100, "最多的彩票数")
var quantum2 = flag.Int("q", 1, "时间片长")
var solve2 = flag.Bool("c", false, "计算答案")

func main() {
	flag.Parse()
	rand.Seed(*seed2)

	tickTotal := 0
	runTotal := 0
	var jobs []*Job2

	if *jlist2 != "" {
		jlistArr := strings.Split(*jlist2, ",")
		*jobnum = len(jlistArr)
		jobs = make([]*Job2, *jobnum)
		for i, jobStr := range jlistArr {
			jobArr := strings.Split(jobStr, ":")
			runTime, _ := strconv.Atoi(jobArr[0])
			tickets, _ := strconv.Atoi(jobArr[1])
			job := &Job2{
				num:     i,
				runTime: runTime,
				tickets: tickets,
			}
			runTotal += runTime
			tickTotal += tickets
			jobs[i] = job
			fmt.Printf("任务 %d ( 运行时间 = %d, 彩票数 = %d )\n", job.num, job.runTime, job.tickets)
		}
	} else {
		jobs = make([]*Job2, *jobnum)
		for i := 0; i < *jobnum; i++ {
			runTime := rand.Intn(*maxlen2) + 1
			tickets := rand.Intn(*maxticket) + 1
			job := &Job2{
				num:     i,
				runTime: runTime,
				tickets: tickets,
			}
			runTotal += runTime
			tickTotal += tickets
			jobs[i] = job
			fmt.Printf("任务 %d ( 运行时间 = %d, 彩票数 = %d )\n", job.num, job.runTime, job.tickets)
		}
	}

	if !*solve2 {
		for i := 0; i < runTotal; i++ {
			r := rand.Intn(1000000) + 1
			fmt.Println("随机数", r)
		}
		return
	}

	var job *Job2
	for i := 0; i < runTotal; {
		r := rand.Intn(1000000) + 1
		winner := r % tickTotal

		current := 0
		for _, jobItem := range jobs {
			current += jobItem.tickets
			if current > winner {
				job = jobItem
				break
			}
		}
		fmt.Printf("随机数 %d -> 抽奖票 %d (总票数 %d) -> 中奖任务 %d\n", r, winner, tickTotal, job.num)
		fmt.Println("\tjobs：")
		for _, jobItem := range jobs {
			wstr := " "
			if jobItem.num == job.num {
				wstr = "*"
			}
			tstr := "---"
			if jobItem.tickets != 0 {
				tstr = strconv.Itoa(jobItem.tickets)
			}
			fmt.Printf("\t(%s 任务:%d 剩余运行时间:%d 票数:%s )\n", wstr, jobItem.num, jobItem.runTime, tstr)
		}
		fmt.Println("")

		if job.runTime > *quantum2 {
			job.runTime -= *quantum2
		} else {
			job.runTime = 0
		}
		i += *quantum2

		if job.runTime == 0 {
			fmt.Printf("--> 任务 %d 完成时间 %d\n", job.num, i)
			tickTotal -= job.tickets
			job.tickets = 0
			*jobnum--
		}

		if *jobnum == 0 {
			return
		}
	}
}
