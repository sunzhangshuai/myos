{
    if (checkComma($2)) {
        gsub(",", "", $2)
        sum += $2
    }
}
END {
    printf("%.5f\n", sum)
}
function checkComma(v) {
    return v ~ /^[+-]?[0-9][0-9]?[0-9]?(,[0-9][0-9][0-9])*([.][0-9]*)?$/
}