{
    gsub(/[.,:;!?(){}]/, "")
    for (i = 1; i <= NF; i++) {
        "echo " $i "| tr a-z A-Z "  | getline key
        close("echo " $i "| tr a-z A-Z ")
        count[key]++
    }
}

END {
    for (w in count) {
        print count[w], w | "sort -rn"
    }
}