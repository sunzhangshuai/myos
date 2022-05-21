# 每一个程序都由一种特殊行开始, 这个特殊行以 .P1 打头,
# 程序以一种特殊行结束, 该行以 .P2 打头。
# 符号对不能嵌套或重叠。
BEGIN {
    expects["aa"] = "bb"
    expects["cc"] = "dd"
    expects["ee"] = "ff"
}

/^(aa|cc|ee)/ {
    if (p != "")
        print "line", NR, ": expected " p
    p = expects[substr($0, 1, 2)]
}

/^(bb|dd|ff)/ {
    x = substr($0, 1, 2)
    if (p != x) {
        print "line", NR, ": saw " x
        if (p)
            print ", expected", p
    }
    p = ""
}
END {
    if (p != "")
        print "at end, missing", p
}