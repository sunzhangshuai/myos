BEGIN {
    RS = "。"
    print "数据太长，只展示5行数据。"
}

NR < 5 {
    gsub(" ", "")
    print length($0), $0 "。"
}

END {
    print "total line：", NR
}