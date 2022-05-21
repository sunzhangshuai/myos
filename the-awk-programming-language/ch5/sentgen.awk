BEGIN {
    # 获取左部出现的次数、每次出现左部拥有右部个数、每次出现的右部权重、以及组成部分
    while ((getline < ARGV[1]) > 0) {
        if ($2 == "->") {
            leftCount[$1]++
            leftRightCount[$1, leftCount[$1]] = NF - 3
            leftRightWidth[$1, leftCount[$1]] = $NF
            for (i = 3; i <= NF - 1; i++) {
                leftRightItem[$1, leftCount[$1], i - 2] = $i
            }
        } else {
            print "illegal production: " $0
        }
    }
    ARGV[1] = ""
    srand()
}

{
    if ($1 in leftCount) {
        print gen($1, 1)
    } else
        print "unknown nonterminal: " $0
}

function gen(sym, vdepth, idx, i, res) {
    if (vdepth > depth) {
        return " " sym
    }
    res = ""
    if (sym in leftCount) {
        idx = randWidthIdx(sym)
        for (i = 1; i <= leftRightCount[sym, idx]; i++) {
            res = res gen(leftRightItem[sym, idx, i], vdepth + 1)
        }
    } else {
        res = " " sym
    }
    return res
}

# 获取带权重的idx
function randWidthIdx(sym, width, count, preNum, i, n) {
    split("", width, "")
    count = leftCount[sym]
    preNum = 0
    for (i = 1; i <= count; i++) {
        preNum += leftRightWidth[sym, i]
        width[preNum] = i
    }
    n = randInt(preNum)
    for (i in width) {
        if (int(n) <= int(i)) {
            return width[i]
        }
    }
    return count
}

# 返回一个 1 - n的随机整数
function randInt(n) {
    return int(n * rand()) + 1
}