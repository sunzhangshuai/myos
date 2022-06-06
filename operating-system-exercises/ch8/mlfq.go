package main

import "flag"

var seed = flag.Int64("s", 0, "指定随机种子")
var numQueues = flag.Int("n", 3, "MLFQ中的队列数（如果没有使用-Q）")
var quantum = flag.Int("q", 10, "时间片长度(如果没有使用-Q参数)")
var allotment = flag.Int("a", 1, "分配长度（如果不使用-A）")
var quantumList = flag.String("Q", "", "指定为x，y，z，...为每个队列级别的时间片长度，其中x是优先级最高的队列的时间片长度，y是第二高的队列的时间片长度，依此类推")
var allotmentList = flag.String("A", "", "x、 y，z，。。。其中x是最高优先级队列，y次高，依此类推")
var numJobs = flag.Int("j", 3, "系统中的作业数")
var maxlen = flag.Int("m", 100, "作业的最大运行时间（如果是随机的）")
var maxio = flag.Int("M", 10, "作业的最大I/O频率（如果是随机的）")
var boost = flag.Int("B", 0, "将所有作业的优先级提高到高优先级的频率（0表示从不）")
var ioTime = flag.Int("i", 5, "I/O 持续时间（固定常数)")
var stay = flag.Bool("S", false, "发出I/O时重置并保持相同的优先级")
var iobump = flag.Bool("I", false, "如果指定，完成I/O的作业将立即移动到当前队列前面")
var jlist = flag.String("l", "", "以逗号分隔的要运行的作业列表，格式为x1，y1，z1：x2，y2，z2：...。其中x是开始时间，y是运行时间，z是作业I/O的频率")
var solve = flag.Bool("c", false, "计算答案")

func main() {
	flag.Parse()

}
