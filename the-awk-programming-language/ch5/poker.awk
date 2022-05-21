BEGIN {
    split(permute(52, 52), pokers, " ")
    # 1. 发牌
    dealPoker(pokers, NORTH, WEST, SOUTH, EAST)
    # 2. 整排
    processPoker(NORTH, north)
    processPoker(WEST, west)
    processPoker(SOUTH, south)
    processPoker(EAST, east)
    # 3. 展示
    show(north, west, south, east)
}

# 洗牌
function permute(k, n,i, p, r) {
    srand()
    p = " "
    for (i = n-k+1; i <= n; i++) {
        r = int(i*rand())+1
        if (p ~ " " r " " ) {
            sub(" " r " ", " " r " " i " ", p)
        } else {
            p = " " r p
        }
    }
    return p
}

# 发牌
function dealPoker(pokers, NORTH, WEST, SOUTH, EAST, i) {
    split("", NORTH, "")
    split("", WEST, "")
    split("", SOUTH, "")
    split("", EAST, "")
    for(i = 1; i <= 52; i++) {
        switch (i % 4) {
            case 0:
                # 北
                handPoker(NORTH, pokers[i])
                break
            case 1:
                # 西
                handPoker(WEST, pokers[i])
                break
            case 2:
                # 南
                handPoker(SOUTH, pokers[i])
                break
            case 3:
                # 东
                handPoker(EAST, pokers[i])
                break
        }
    }
}

# 拿牌
function handPoker(p, v, i, pos, len) {
    len = length(p)
    pos = len + 1
    for (i = 1; i <= len; i++) {
        if (p[i] < v) {
            pos = i
            break
        }
    }
    if (pos <= len) {
        for (i = len + 1; i > pos; i--) {
            p[i] = p[i - 1]
        }
    }
    p[pos] = v
}

# 整牌
function processPoker(p, res, i) {
    res["S"] = "S:"
    res["H"] = "H:"
    res["D"] = "D:"
    res["C"] = "C:"
    for (i in p) {
        switch (int((p[i] - 1) / 13)) {
            case 0:
                res["S"] = res["S"] " " fvcard(p[i])
                break
            case 1:
                res["H"] = res["H"] " " fvcard(p[i])
                break
            case 2:
                res["D"] = res["D"] " " fvcard(p[i])
                break
            case 3:
                res["C"] = res["C"] " " fvcard(p[i])
                break
        }
    }
}

# 输出实际牌号
function fvcard(i) {
    if (i % 13 == 0) return "A"
    else if (i % 13 == 12) return "K"
    else if (i % 13 == 11) return "Q"
    else if (i % 13 == 10) return "J"
    else return (i % 13) + 1
}

# 展示牌
function show(north, west, south, east) {
    printf("%-40s%-40s%-40s\n", "", north["S"],"")
    printf("%-40s%-40s%-40s\n", "", north["H"],"")
    printf("%-40s%-40s%-40s\n", "", north["D"],"")
    printf("%-40s%-40s%-40s\n", "", north["C"],"")
    printf("%-40s%-40s%-40s\n", west["S"], "",east["S"])
    printf("%-40s%-40s%-40s\n", west["H"], "",east["H"])
    printf("%-40s%-40s%-40s\n", west["D"], "",east["D"])
    printf("%-40s%-40s%-40s\n", west["C"], "",east["C"])
    printf("%-40s%-40s%-40s\n", "", south["S"],"")
    printf("%-40s%-40s%-40s\n", "", south["H"],"")
    printf("%-40s%-40s%-40s\n", "", south["D"],"")
    printf("%-40s%-40s%-40s\n", "", south["C"],"")
}