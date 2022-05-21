# 将字段字割符设置为制表符 (\t), 并在输出之前打印标题。
# 每一列都刚 好与标题的列表头对齐。
# 打印总和。
BEGIN {
    FS = "\t"
    printf("%10s %6s %5s %20s\n\n", "country", "area", "pop", "continent")
}
{
    printf("%10s %6s %5s %20s\n", $1, $2, $3, $4)
    area += $2
    pop += $3
}
END {
    printf("\n %15s：%6s %15s：%5s\n", "total area", area, "total pop", pop)
    print FNR
}
