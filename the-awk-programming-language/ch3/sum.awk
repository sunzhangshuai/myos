NR == 1 { maxfield = NF } # 相信程序每行的字段数都一样
NF != 0 {
    for (i = 1; i <= NF; i++) {
        # 判断是否为整数
        if (isnum($i)) {
            sum[i] += $i
        }
    }
}
END {
    for (i = 1; i <= maxfield; i++) {
        if (isnum(sum[i])) {
            printf("%g", sum[i])
        } else {
            printf("%s", "--")
        }
        printf("%s", i == maxfield ? "\n" : " ")
    }
}
function isnum(v) {
    return v ~ /^(\+|-)?[0-9]+\.?[0-9]*$/
}