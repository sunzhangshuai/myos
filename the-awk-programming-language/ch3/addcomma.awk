BEGIN {
    system("cat /dev/null > ch3/file/sumcommatest.txt")
}
{
    printf("%-12s %20s\n", $0, addcomma($0)) >> "ch3/file/sumcommatest.txt"
}
function addcomma(v, num) {
    if (v < 0) {
        return "-" addcomma(-v)
    }
    num = sprintf("%.2f", v)
    while (num ~ /[0-9][0-9][0-9][0-9]/) {
        sub(/[0-9][0-9][0-9][,.]/, ",&", num)
    }
    return num
}