# 预处理&准备与排序
BEGIN { FS = "\t" }
{
    den = 1000 * $3/$2
    printf("%-15s:%12.8f:%s:%d:%d:%.1f\n", $4, 1/den, $1, $3, $2, den) | "sort"
}