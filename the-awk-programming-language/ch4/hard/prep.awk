BEGIN  {
    FS = "\t"
}
pass == 1 {
    areatot += $2
    poptot += $3
}

pass == 2 {
    den = 1000*$3/$2
    printf("%s:%s:%s:%f:%d:%f:%f\n", $4, $1, $3, 100*$3/poptot, $2, 100*$2/areatot, den) | "sort -t: +0 -1 +6rn"
}