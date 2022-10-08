package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var seed7 = flag.Int64("s", 0, "指定随机种子")
var addressesstr = flag.String("A", "-1", "要访问的一组逗号分隔的页面；-1表示随机生成")
var asizestr2 = flag.String("a", "16k", "地址空间大小(e.g., 16, 64k, ...)")
var psizestr2 = flag.String("p", "64k", "物理空间大小 (e.g., 16, 64k, ...)")
var pagesizestr = flag.String("P", "4k", "页的大小(e.g., 4k, 8k, ...)")
var num2 = flag.Int("n", 5, "生成的虚拟地址数量")
var usedPercent = flag.Int("u", 50, "已经使用的百分比")
var verbose = flag.Bool("v", false, "详细模式")
var solve7 = flag.Bool("c", false, "计算答案")

func main() {
	flag.Parse()

	rand.Seed(*seed7)
	asize := convert2(*asizestr2)
	psize := convert2(*psizestr2)
	pagesize := convert2(*pagesizestr)

	if psize >= convert2("1g") || asize >= convert2("1g") {
		fmt.Println("错误：此模拟必须使用较小的大小（小于1 GB）。")
		os.Exit(1)
	}

	mustbemultipleof(asize, pagesize, "虚拟空间必须是分页大小的倍数")
	mustbemultipleof(psize, pagesize, "物理空间必须是虚拟空间的倍数")

	// 总页数和虚拟页数
	pages := psize / pagesize
	vpages := asize / pagesize
	if psize < asize {
		fmt.Println("错误：物理内存大小必须大于地址空间大小（对于此模拟）")
		os.Exit(1)
	}

	// 页占用情况
	used := make([]int, pages)
	// 虚拟页对应的物理页
	pt := make(map[int]int)

	vabits := int(math.Log2(float64(asize)))
	mustbepowerof2(vabits, asize, "虚拟空间大小必须是2的幂")
	pagebits := int(math.Log2(float64(pagesize)))
	mustbepowerof2(pagebits, pagesize, "页的大小必须是2的幂")
	// page 偏移掩码，用来表示偏移的位全为1。
	pagemask := (1 << pagebits) - 1
	// vpn 掩码，用来表示vpn的位全为1。
	vpnmask := 0xFFFFFFFF & ^pagemask

	fmt.Println("页面表的格式很简单：")
	fmt.Println("高阶（最左侧）位是有效位。")
	fmt.Println("\t如果位为1，则条目的其余部分为PFN。")
	fmt.Println("\t如果位为0，则页面无效。")
	fmt.Println("如果要按页面表的每个条目打印VPN，请使用详细模式（-v）。")
	fmt.Println()

	fmt.Println("页表（从条目0到最大）")

	// 初始化已有页面，有一定的无效概率
	for i := 0; i < vpages; i++ {
		for done := 0; done == 0; {
			// 有效
			if rand.Intn(100) > (100 - *usedPercent) {
				u := rand.Intn(pages)
				if used[u] == 0 {
					used[u] = 1
					done = 1
					// 设置位有效
					if *verbose {
						fmt.Printf("  [%8d]  0x%08x\n", i, 0x80000000|u)
					} else {
						fmt.Printf("  0x%08x\n", 0x80000000|u)
					}
					pt[i] = u
				}
			} else {
				// 设置位无效
				if *verbose {
					fmt.Printf("  [%8d]  0x%08x\n", i, 0)
				} else {
					fmt.Printf("  0x%08x\n", 0)
				}
				pt[i] = -1
				done = 1
			}
			done = 1
		}
	}
	fmt.Println()

	// 获取地址列表
	var addresses []int
	if *addressesstr == "-1" {
		addresses = make([]int, *num2)
		for i := 0; i < *num2; i++ {
			addresses[i] = rand.Intn(asize)
		}
	} else {
		addrList := strings.Split(*addressesstr, ",")
		*num2 = len(addrList)
		addresses = make([]int, *num2)
		for i := 0; i < *num2; i++ {
			addresses[i], _ = strconv.Atoi(addrList[i])
		}
	}

	fmt.Println("虚拟地址调用栈：")
	if !*solve7 {
		for _, address := range addresses {
			fmt.Printf("  虚拟地址 0x%08x (decimal: %8d) --> 物理地址 或 无效地址 ？\n", address, address)
		}

		fmt.Println("对于每个虚拟地址，记下它转换为的物理地址，或记下它是一个越界地址（例如segfault）。")
		fmt.Println()
		return
	}

	for _, address := range addresses {
		paddr := 0
		// 页，拿到地址表示页的数值
		vpn := (address & vpnmask) >> pagebits
		pfn, ok := pt[vpn]
		if pfn < 0 || !ok {
			fmt.Printf("  虚拟地址 0x%08x (decimal: %8d) -->  无效 (虚拟页 %d 无效)\n", address, address, vpn)
		} else {
			offset := address & pagemask // 获取偏移量
			paddr = (pfn << pagebits) | offset
			fmt.Printf("  虚拟地址 0x%08x (decimal: %8d) -->  %08x (decimal %8d)[VPN %d]\n", address, address, paddr, paddr, vpn)
		}
	}
}

// convert2 大小转化
func convert2(size string) int {
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

// mustbemultipleof 必须是倍数
func mustbemultipleof(bignum int, num int, msg string) {
	if bignum%num != 0 {
		fmt.Println("参数错误：", msg)
		os.Exit(1)
	}
}

// mustbepowerof2 size 必须是 2 的 bit 次方
func mustbepowerof2(bits int, size int, msg string) {
	if math.Pow(2, float64(bits)) != float64(size) {
		fmt.Println("参数错误：", msg)
		os.Exit(1)
	}
}
