BEGIN  {
    FS = ":"
    blanks = blanks_ = blanks__ = sprintf("%100s", " ")
    gsub(/ /, "-", blanks_)
    gsub(/ /, "=", blanks__)
}

{
    row[NR] = $0
    for (i = 1; i <= 7; i++) {
        fieldWidth[i] = max(fieldWidth[i], length($i) + 6)
    }
    gpop += $3
    gpoppct += $4
    garea += $5
    gareapct += $6
}

END {
    # 报表编号  人口 地区 人口密度  日期
    printf("%s%s%s\n\n",
        outputField("Report No. 1", fieldWidth[1]),
        outputField("POPULATION, AREA, POPULATION DENSITY", fieldWidth[2] + fieldWidth[3] + fieldWidth[4] + fieldWidth[5]),
        outputField(getDate(), fieldWidth[6] + fieldWidth[7]))

    # 大陆 国家       人口             面积             人口密度
    printf("%s%s%s%s%s\n\n",
        outputField("CONTINENT", fieldWidth[1]),
        outputField("COUNTRY", fieldWidth[2]),
        outputField("POPULATION", fieldWidth[3] + fieldWidth[4]),
        outputField("AREA", fieldWidth[5] + fieldWidth[6]),
        outputField("POP. DEN.", fieldWidth[7]))

    #           百万人 百分比     平方千米 百分比    人口密度【每平方米】
    printf("%s%s%s%s%s%s%s\n",
        outputField("", fieldWidth[1]),
        outputField("", fieldWidth[2]),
        outputField("Millions", fieldWidth[3]),
        outputField("Pct. of", fieldWidth[4]),
        outputField("Thousands", fieldWidth[5]),
        outputField("Pct. of", fieldWidth[6]),
        outputField("People per", fieldWidth[7]))
    printf("%s%s%s%s%s%s%s\n",
        outputField("", fieldWidth[1]),
        outputField("", fieldWidth[2]),
        outputField("of People", fieldWidth[3]),
        outputField("Total", fieldWidth[4]),
        outputField("of Sq. Mi.", fieldWidth[5]),
        outputField("Total", fieldWidth[6]),
        outputField("Sq. Mi.", fieldWidth[7]))
    printf("%s%s\n",
        outputField("", fieldWidth[1] + fieldWidth[2]),
        outputFieldByBlank("", fieldWidth[3] + fieldWidth[4] + fieldWidth[5] + fieldWidth[6], blanks_))

    for (r = 1; r <= NR; r ++) {
        split(row[r], d, ":")
        if (d[1] != prev) {
            if (r > 1) {
                totalprint()
            }
            prev = d[1]
            poptot = d[3]
            poppct = d[4]
            areatot = d[5]
            areapct = d[6]
        } else {
            d[1] = ""
            poptot += d[3]
            poppct += d[4]
            areatot += d[5]
            areapct += d[6]
        }
        printf("%s%s%s%s%s%s%s\n",
            outputField(d[1], fieldWidth[1]),
            outputField(d[2], fieldWidth[2]),
            outputField(d[3], fieldWidth[3]),
            outputField(d[4], fieldWidth[4]),
            outputField(d[5], fieldWidth[5]),
            outputField(d[6], fieldWidth[6]),
            outputField(d[7], fieldWidth[7]))
    }

    totalprint()
    printf("%s%s%s%s%s\n",
        outputField("GRAND TOTAL", fieldWidth[1] + fieldWidth[2]),
        outputField(gpop, fieldWidth[3]),
        outputField(gpoppct, fieldWidth[4]),
        outputField(garea, fieldWidth[5]),
        outputField(gareapct, fieldWidth[6]))
    printf("%s%s\n",
        outputField("", fieldWidth[1] + fieldWidth[2]),
        outputFieldByBlank("", fieldWidth[3] + fieldWidth[4] + fieldWidth[5] + fieldWidth[6], blanks__))
}

# 打印大陆总数
function totalprint() {     # print totals for previous continent
    printf("%s%s\n",
        outputField("", fieldWidth[1] + fieldWidth[2]),
        outputFieldByBlank("", fieldWidth[3] + fieldWidth[4] + fieldWidth[5] + fieldWidth[6], blanks_))
    printf("%s%s%s%s%s\n",
        outputField("TOTAL for " prev, fieldWidth[1] + fieldWidth[2]),
        outputField(poptot, fieldWidth[3]),
        outputField(poppct, fieldWidth[4]),
        outputField(areatot, fieldWidth[5]),
        outputField(areapct, fieldWidth[6]))
    printf("%s%s\n",
        outputField("", fieldWidth[1] + fieldWidth[2]),
        outputFieldByBlank("", fieldWidth[3] + fieldWidth[4] + fieldWidth[5] + fieldWidth[6], blanks__))

}

function max(x, y) { return (x > y) ? x : y }

function outputFieldByBlank(fieldValue, fieldWidth, blanks, rigth) {
    left = int((fieldWidth - length(fieldValue)) / 2)
    return substr(blanks, 1, fieldWidth - length(fieldValue) - left) fieldValue substr(blanks, 1, left)
}

function outputField(fieldValue, fieldWidth) {
    return outputFieldByBlank(fieldValue, fieldWidth, blanks)
}

function getDate(date, d) {
    "date" | getline date
    split(date, d, " ")
    date = d[2] " " d[3] ", " d[6]
    return date
}