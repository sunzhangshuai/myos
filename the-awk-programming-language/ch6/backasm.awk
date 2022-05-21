BEGIN {
    srcfile = ARGV[1]
    ARGV[1] = ""
    n = split("const get put ld st add sub jpos jz j halt", x)

    # 第一遍获取所有操作符和引用
    nextmem = 0
    while (getline < srcfile > 0) {
        op[nextmem] = x[int(substr($1, 1, 2)) + 1]
        quote[nextmem] = int(substr($1, 3))
        nextmem++
    }
    close(srcfile)

    # 第二遍判断引用位置是否需要符号
    constIndex = symbolIndex = 1
    for (i = 0; i < nextmem; i++) {
        if (!(quote[i] in symbol) && quote[i] != 0) {
            if (op[quote[i]] == "const") {
                symbol[quote[i]] = "const" constIndex++
            } else {
                symbol[quote[i]] = "symbol" constIndex++
            }
        }
    }

    # 第三遍生成汇编文件
    for (i = 0; i < nextmem; i++) {
        # 先打印符号
        if (i in symbol) {
            printf("%s", symbol[i])
        }
        # 打印操作
        printf("\t%s", op[i])
        # 打印操作
        if (op[i] == "const") {
            printf("\t%s\n", quote[i])
            continue
        } else {
            if (quote[i] != 0) {
                printf("\t%s", symbol[quote[i]])
            }
        }
        printf("\n")
    }
}