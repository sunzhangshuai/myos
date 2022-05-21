{
    # 画布长宽
    ht = 50
    wid = 150

    # 原点偏移
    ox = 6
    oy = 2

    # 数字验证规则
    number = "^[-+]?([0-9]+[.]?[0-9]*|[.][0-9]+)([eE][-+]?[0-9]+)?$"

    nb = nl = 10
}

# 获取图表描述，展示在底部
$1 == "label" {
    sub(/^ *label */, "")
    botlab = $0
    next
}

## 底部横轴标记
#$1 == "bottom" && $2 == "ticks" {
#    for (i = 3; i <= NF; i++) {
#        bticks[++nb] = $i
#    }
#    next
#}
#
## 左侧纵轴标记
#$1 == "left" && $2 == "ticks" {
#    for (i = 3; i <= NF; i++) {
#        lticks[++nl] = $i
#    }
#    next
#}

# 图表纵轴和横轴的最小最大值
#$1 == "range" {
#    xmin = $2
#    ymin = $3
#    xmax = $4
#    ymax = $5
#    next
#}

# 如果指定宽高的话覆盖默认宽高
$1 == "height" {
    ht = $2
    next
}
$1=="width" {
    wid=$2
    next
}

# 点坐标
$1 ~ number && $2 ~ number {
    nd++
    if (xmax == 0) {
        xmax = xmin = $1
        ymax = ymin = $2
    } else {
        if (xmax < $1) xmax = $1
        if (xmin > $1) xmin = $1
        if (ymax < $2) ymax = $2
        if (ymin > $2) ymin = $2
    }
    x[nd] = $1
    y[nd] = $2
    next
}

# 画图
END {
    # 初始化标记
    initTicks()
    # 画框框
    frame()
    # 画标点
    ticks()
    # 画底部label
    label()
    # 画数据
    data()
    # draw
    draw()
}

# 初始化标记
function initTicks() {
    xrange = xmax - xmin
    for (i = 0; i <= nb - 1; i++) {
        bticks[i+1] = int(xrange / (nl - 1) * i + xmin + 0.5)
    }
    xmax += int(xrange / 20)
    xmin -= int(xrange / 20)
    yrange = ymax - ymin
    for (i = 0; i <= nl - 1; i++) {
        lticks[i+1] = int(yrange / (nl - 1) * i + ymin + 0.5)
    }
    ymax += int(yrange / 20)
    ymin -= int(yrange / 20)
}

# 画框框
function frame() {
    for (i = ox; i < wid; i++) {
        plot(i, oy, "-")
        plot(i, ht - 1, "-")
    }
    for (i = oy; i < ht; i++) {
        plot(ox, i, "|")
        plot(wid - 1, i, "|")
    }
}

# 画标记
function ticks() {
    for (i = 1; i <= nb; i++) {
        plot(xscale(bticks[i]), oy, "|")
        splot(xscale(bticks[i]) - 1, 1, bticks[i])
    }
    for (i = 1; i <= nl; i++) {
        plot(ox, yscale(lticks[i]), "-")
        splot(0, yscale(lticks[i]), lticks[i])
    }
}

# 画底部label
function label() {
    splot(int((wid + ox - length(botlab))/2), 0, botlab)
}

# 画数据
function data() {
    for (i = 1; i <= nd; i++) {
        plot(xscale(x[i]), yscale(y[i]), "*")
    }
}

# 画
function draw(   i, j) {
    for (i = ht - 1; i >= 0; i--) {
        for (j = 0; j < wid; j++) {
            printf((j,i) in array ? array[j,i] : " ")
        }
        printf("\n")
    }
}

# 布局字符
function plot(x, y, c) {
    array[x, y] = c
}

# 布局字符串
function splot(x, y, s,    i, n) {
    n = length(s)
    for (i = 0; i < n; i++) {
        array[x+i, y] = substr(s, i+1, 1)
    }
}

# 计算横轴相对值
function xscale(i) {
    return int((i - xmin) / (xmax - xmin) * (wid - 1 - ox) + ox + 0.5)
}

# 计算纵轴相对值
function yscale(i) {
    return int((i - ymin) / (ymax - ymin) * (ht - 1 - oy) + oy + 0.5)
}