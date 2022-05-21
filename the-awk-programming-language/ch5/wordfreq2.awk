{
    if (prev == "") {
        prev = $0
    }
    if (prev != $0) {
        print prev, num | "sort  -k2rn"
        num = 0
        prev = $0
    } else {
        num ++
    }
}

END {
    print prev, num | "sort -k2rn"
}