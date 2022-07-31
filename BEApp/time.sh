t1=`date +'%Y-%m-%d %H:%M:%S'`
start_seconds=$(date --date="$t1" +%s);

echo "$t1: 实验开始："

./main

t2=`date +'%Y-%m-%d %H:%M:%S'` 
echo "$t2: 第$i次实验结束"
end_seconds=$(date --date="$t2" +%s);

echo "运行时间： "$((end_seconds-start_seconds))"s"