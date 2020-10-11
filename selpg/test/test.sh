selpg -s2 -e3 testfile >output1
diff output1 ans1

selpg -s10 -e20 -l10 testfile >output2
diff output2 ans2

selpg -s10 -e20 -l10 <testfile >output2
diff output2 ans2

selpg -s10 -e20 -l10 testfile | cat >output2
diff output2 ans2

cat testfile | selpg -s10 -e20 -l10 >output2
diff output2 ans2

selpg -s10 -e20 -l10 null 2>output3
diff output3 ans3

selpg -s10 -e20 -l10 testfile >output2 2>output3
diff output2 ans2

selpg -s10 -e20 -l10 testfile >output2 2>/dev/null
diff output2 ans2

selpg -s10 -e20 -l10 testfile >/dev/null
diff output2 ans2

selpg -s3 -e4 -f testfile >output4
diff output4 ans4