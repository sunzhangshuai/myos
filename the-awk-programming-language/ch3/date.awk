# 计算到1970-01-01的间隔天数
function datenum(year, month, day, days, i, n) {
    split("31 28 31 30 31 30 31 31 30 31 30 31", days)
    n = (year - 1901) * 365 + int((year - 1901) / 4)
    if (year % 4 == 0) {
        days[2]++
    }
    for (i = 1; i < month; i++) {
        n += days[i]
    }
    return n + day
}
{
    print datenum(substr($1, 5, 2) + 1900, substr($1, 1, 2), substr($1, 3, 2)), $2
}