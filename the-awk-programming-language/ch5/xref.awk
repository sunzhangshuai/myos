/^\.#/ {
    printf("{ gsub(/%s/, \"%d\") }\n", $2, ++count[$1])
}
END    {
    printf("!/^[.]#/\n")
}