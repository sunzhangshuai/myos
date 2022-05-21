{
    # 数字验证规则
    number = "^[-+]?([0-9]+[.]?[0-9]*|[.][0-9]+)([eE][-+]?[0-9]+)?$"
}

# 获取图表描述，展示在底部
$1 == "label" {
    print $0
}

# 底部横轴标记
$1 == "bottom" && $2 == "ticks" {
    $1 = "left"
    print
    next
}

# 左侧纵轴标记
$1 == "left" && $2 == "ticks" {
    $1 = "bottom"
    print
    next
}

# 图表纵轴和横轴的最小最大值
$1 == "range" {
    print "range", $3, $2, $5, $4
    next
}

# 如果指定宽高的话覆盖默认宽高
$1 == "height" {
    $1 = "width"
    print
    next
}
$1=="width" {
    $1 = "height"
    print
    next
}

# 点坐标
$1 ~ number && $2 ~ number {
    print $2, $1
    next
}