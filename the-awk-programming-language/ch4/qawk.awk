BEGIN {
    # 初始化关系
    readrel("file/relfile.txt")
    RS = ""
}

/./ {
    # 清理历史查询数据
    for (i in qattr) {
        delete qattr[i]
    }
    print "origin query：" $0
    doquery($0)
    print ""
}

# 最终初始化出一下数据，用于保存关系数据
# relname 表的列表，nrel表的个数
# cmd 表的命令，ncmd为表的命令计数，ncmd 每张表的命令数
# attr 每张表、每个字段的位置，nattr 每张表的字段数
function readrel(f) {
    while ((getline <f) > 0)
        if ($0 ~ /^[A-Za-z\/\.]+ *:/) {
            # 拿到表名
            gsub(/[^A-Za-z\/\.]+/, "", $0)
            relname[++nrel] = $0
        } else if ($0 ~ /^[ \t]*!/)     # !命令
            cmd[nrel, ++ncmd[nrel]] = substr($0,index($0,"!")+1)
        else if ($0 ~ /^[ \t]*[A-Za-z]+[ \t]*$/)  # attribute
            attr[nrel, $1] = ++nattr[nrel]
        else if ($0 !~ /^[ \t]*$/)      # not white space
            print "bad line in relfile:", $0
}

# 执行查询
function doquery(s, i,j) {
    # 拿到原始query
    query = s

    # 判断是否还有变量可以获取
    while (match(s, /\$[A-Za-z]+/)) {
        # 提取变量
        qattr[substr(s, RSTART+1, RLENGTH-1)] = 1
        s = substr(s, RSTART+RLENGTH+1)
    }

    # 查询哪张表能包含所有的字段，i值即为表的索引
    for (i = 1; i <= nrel && !subset(qattr, attr, i); ) {
        i++
    }

    # 没有表能包含所有字段
    if (i > nrel) {
        missing(qattr)
    } else {
        # 去掉所有的 $ 符号
        for (j in qattr) {
            gsub("\\$" j, "$" attr[i,j], query)
        }

        systemCmd = ""
        if (exists[i] && ncmd[i] > 0) {
            # 如果当前表需要命令生成，则执行命令
            for (j = 1; j <= ncmd[i]; j++) {
                systemCmd = systemCmd cmd[i, j] "\n"
            }
            if (system(systemCmd) != 0) {
                print "command failed, query skipped\n", systemCmd
                return
            }
            exists[i]++
        }

        # 对表文件执行 awk
        awkcmd = sprintf("awk -F'\\t' '%s' %s", query, relname[i])
        printf("query: %s\n", awkcmd)
        system(awkcmd)
    }

}

# 判断
function subset(q, a, r, i) {
    for (i in q) {
        if (!((r,i) in a)) {
            return 0
        }
    }
    return 1
}
function missing(x, i) {
    print "no table contains all of the following attributes:"
    for (i in x)
        print i
}

