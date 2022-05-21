BEGIN {
    blanks = sprintf("%60s", " ")
}

/./  {
    for (i = 1; i <= NF; i++) {
        addword($i)
    }
}

/^$/ {
    printline("no"); print ""
}

END  {
    printline("no")
}

# 添加词
function addword(w) {
    if (cnt + size + length(w) > 60)
        printline()
    line[++cnt] = w
    size += length(w)
}

# 打印当前行
function printline(f, i, nsp, spaceNum, holeNum) {
    # 参数 no 可以避免对段落的最后一行进行右对齐
    if (f == "no" || cnt == 1) {
        for (i = 1; i <= cnt; i++) {
            printf("%s%s", line[i], i < cnt ? " " : "\n")
        }
    } else if (cnt > 1) {
        # 空格的间隙交错开，一次向上取整，一次向下取整
        dir = 1 - dir
        spaceNum = 60 - size
        holeNum = cnt - 1
        for (i = 1; i <= cnt - 1; i++) {
            nsp = int((spaceNum -dir) / holeNum) + dir
            printf("%s%s", line[i], substr(blanks, 1, nsp))
            holeNum--
            spaceNum -= nsp
        }
        print line[cnt]
    }
    cnt = size = 0
}