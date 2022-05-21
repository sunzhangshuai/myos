BEGIN {
    key = 0
}

/ no |not |n't / { print "error: can't do negatives:", $0; ok = 0 }

# 必须匹配到一个规则
{ ok = 0 }

# 去重
/uniq|discard.*(iden|dupl)/ { uniq = " -u"; ok = 1 }

# 分隔符
/separ.*tab|tab.*sep/ { sep = "t'\t'"; ok = 1 }
/separ/ {
    # 只拿只有一个字符的字段
    for (i = 1; i <= NF; i++) {
        if (length($i) == 1) {
            sep = "t'" $i "'"
        }
    }
    ok = 1
}

# 排序项
/key/ {
    key++; dokey();ok = 1
}

# 排序方式
/dict/                            { dict[key] = "d"; ok = 1 }
/ignore.*(space|blank)/           { blank[key] = "b"; ok = 1 }
/fold|case/                       { fold[key] = "f"; ok = 1 }
/num/                             { num[key] = "n"; ok = 1 }
/rev|descend|decreas|down|oppos/  { rev[key] = "r"; ok = 1 }
/forward|ascend|increas|up|alpha/ { ok = 1 }

# 匹配不到报错
!ok   { print "error: can't understand:", $0 }

END {
    cmd = "sort" uniq
    flag = dict[0] blank[0] fold[0] rev[0] num[0] sep
    if (flag) cmd = cmd " -" flag
    for (i = 1; i <= key; i++) {
        if (pos[i] != "") {
            flag = pos[i] dict[i] blank[i] fold[i] rev[i] num[i]
            if (flag) cmd = cmd " +" flag
            if (pos2[i]) cmd = cmd " -" pos2[i]
        }
    }
    print cmd
}

function dokey(   i) {
    for (i = 1; i <= NF; i++) {
        if ($i ~ /^[0-9]+$/) {
            pos[key] = $i - 1
            break
        }
    }
    for (i++; i <= NF; i++) {
        if ($i ~ /^[0-9]+$/) {
            pos2[key] = $i
            break
        }
    }
    if (pos[key] == "") {
        printf("error: invalid key specification: %s\n", $0)
    }
    if (pos2[key] == "") {
        pos2[key] = pos[key] + 1
    }
}