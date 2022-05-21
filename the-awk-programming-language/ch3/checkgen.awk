BEGIN {
    FS = "\t+"
}
$1 == "BEGIN" {
    printf("BEGIN {\n\t%s\n}\n", $2)
    next
}
{
    printf("%s {\n\tprintf(\"line %%d, %s: %%s\\n\",NR,$0)\n}\n",$2, $3)
}