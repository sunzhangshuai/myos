BEGIN {
    srcfile = ARGV[1]
    ARGV[1] = ""
    tempfile = "ch6/file/asm.temp"
    backfile = "ch6/file/backasm.txt"
    n = split("const get put ld st add sub jpos jz j halt", x)
    for (i = 1; i <= n; i++) {
        # 生成每个指令的操作码
        op[x[i]] = i - 1
    }

    # 汇编第一步
    nextmem = 0
    FS = "[ \t]+"
    while (getline < srcfile > 0) {
        input[nextmem] = $0
        # 清空注释
        sub(/#.*/, "")
        # 统计标签及地址
        symtab[$1] = nextmem
        if ($2 != "") {
            print $2 "\t" $3 > tempfile
            nextmem++
        }
    }
    close(tempfile)

    # 汇编第二步
    nextmem = 0
    while (getline <tempfile > 0) {
        # 如果是符号地址的话
        if ($2 !~ /^[0-9]*$/) {
            # 将符号替换成地址
            $2 = symtab[$2]
        }
        # 保存每行指令
        mem[nextmem++] = 1000 * op[$1] + $2  # pack into word
    }
    for (i = 0; i < nextmem; i++) {
        printf("%3d:  %05d   %s\n", i, mem[i], input[i])
        printf("%3d:  %05d   %s\n", i, mem[i], input[i]) > backfile
    }
    print "=============================================================="
    close(backfile)

    if (ARGV[2] == "auto") {
        exit
    }
    # 执行
    for (pc = 0; pc >= 0;) {
        addr = mem[pc] % 1000
        code = int(mem[pc++] / 1000)
        sub(/#.*/, "", input[pc])
        sub(/ *$/, "", input[pc])
        printf("const sum=%d; %s: \n", mem[10], mem[11], input[pc - 1])
        if      (code == op["get"])  { getline acc }
        else if (code == op["put"])  { print acc }
        else if (code == op["st"])   { mem[addr] = acc }
        else if (code == op["ld"])   { acc  = mem[addr] }
        else if (code == op["add"])  { acc += mem[addr] }
        else if (code == op["sub"])  { acc -= mem[addr] }
        else if (code == op["jpos"]) { if (acc >  0) pc = addr }
        else if (code == op["jz"])   { if (acc == 0) pc = addr }
        else if (code == op["j"])    { pc = addr }
        else if (code == op["halt"]) { pc = -1 }
        else                         { pc = -1 }
    }
}