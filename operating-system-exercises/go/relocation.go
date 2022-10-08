package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
)

var seed4 = flag.Int64("s", 0, "指定随机种子")
var asizestr = flag.String("a", "1k", "地址空间大小 (e.g., 16, 64k, 32m)")
var psizestr = flag.String("p", "16k", "物理内存大小(e.g., 16, 64k)")
var num = flag.Int("n", 5, "生成虚拟地址数量")
var basestr = flag.String("b", "-1", "基址寄存器的值")
var limitstr = flag.String("l", "-1", "界限寄存器的值")
var solve4 = flag.Bool("c", false, "计算答案")

func main() {
	flag.Parse()

	fmt.Println()
	fmt.Println("参数 随机数种子", *seed4)
	fmt.Println("参数 虚拟地址大小", *asizestr)
	fmt.Println("参数 物理内存大小", *psizestr)
	fmt.Println()

	rand.Seed(*seed4)
	asize := convert(*asizestr)
	psize := convert(*psizestr)
	base := convert(*basestr)
	limit := convert(*limitstr)

	if limit == -1 {
		limit = asize/4 + rand.Intn(asize/4)
	}

	if base == -1 {
		for done := 0; done == 0; done++ {
			base = rand.Intn(psize)
			if base+limit < psize {
				done = 1
			}
		}
	}

	fmt.Println("基址和边界寄存器信息：")
	fmt.Println()
	fmt.Printf("\t基址: 0x%o (decimal %d)\n", base, base)
	fmt.Printf("\t边界: %d", limit)
	fmt.Println()

	fmt.Println("虚拟地址调用栈")
	for i := 0; i < *num; i++ {
		vaddr := rand.Intn(asize)
		if !*solve4 {
			fmt.Printf("虚拟地址%2d: 0x%o (decimal: %4d) --> PA或段冲突\n", i, vaddr, vaddr)
		} else {
			paddr := 0
			if vaddr >= limit {
				fmt.Printf("虚拟地址%2d: 0x%o (decimal: %4d) --> 不合法\n", i, vaddr, vaddr)
			} else {
				paddr = vaddr + base
				fmt.Printf("虚拟地址%2d: 0x%o (decimal: %4d) --> 合法: 0x%o (decimal: %4d)\n", i, vaddr, vaddr, paddr, paddr)
			}
		}
	}

	if !*solve4 {
		fmt.Println("对于每个虚拟地址，写下它转换为的物理地址，或者写下它是一个越界地址（分段冲突）。对于这个问题，您应该假设一个给定大小的简单虚拟地址空间。")
		fmt.Println()
		return
	}
}

// convert 大小转化
func convert(size string) int {
	var res int
	length := len(size)
	lastchar := size[length-1]
	switch lastchar {
	case 'k', 'K':
		res, _ = strconv.Atoi(size[:length-1])
		res *= 1024
	case 'm', 'M':
		res, _ = strconv.Atoi(size[:length-1])
		res *= 1024 * 1024
	case 'g', 'G':
		res, _ = strconv.Atoi(size[:length-1])
		res *= 1024 * 1024 * 1024
	default:
		res, _ = strconv.Atoi(size)
	}
	return res
}
