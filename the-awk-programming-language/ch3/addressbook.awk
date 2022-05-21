# 若记录由空百行分隔，则将RS值设置为空白符
BEGIN {
    RS = ""
    ORS = "\n\n"
    FS = "\n"
}
/New York/ {
    print $1, $4
}