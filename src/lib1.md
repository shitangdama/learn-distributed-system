time go run wc.go master sequential pg-*.txt
sort -n -k2 mrtmp.wcseq | tail -10
rm mrtmp.*