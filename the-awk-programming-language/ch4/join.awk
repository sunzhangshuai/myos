BEGIN {
    OFS = sep = "\t"

    eofstat = 1  # 第二个文件是否读完的标识
    # 存放每组数据 gp
    ungot # 判断是否有缓存数据
    ungotline # 缓存的数据内容
    filename = ARGV[2] # 第二个文件
    ARGV[2] = ""
    preGroup # 前一组标识
    ng = getGroup()
    if (ng <= 0) {
        exit
    }
}

{
    while (prefix($0) > prefix(gp[1])) {
        ng = getGroup()
        if (ng <= 0) {
            exit
        }
    }
    if (prefix($0) == prefix(gp[1]))  {
        for (i = 1; i <= ng; i++) {
            print $0, suffix(gp[i])
        }
    }
}

function getGroup(n) {
    n = 0
    while (eofstat) {
        n++
        getOne(n)
        if (eofstat && (prefix(gp[n]) != prefix(gp[1]))) {
            unget(gp[n])
            break
        }
    }
    return n-1
}

function getOne(n) {
    if (eofstat <= 0) {
        return 0
    }
    if (ungot) {
        gp[n] = ungotline
        ungot = 0
        return 1
    }
    eofstat = getline gp[n] < filename
}

function unget(s)  {
    ungotline = s; ungot = 1
}

function prefix(s) {
    return substr(s, 1, index(s, sep) - 1)
}

function suffix(s) {
    return substr(s, index(s, sep) + 1)
}