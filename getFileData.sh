#!/bin/bash

rm -r /home/k8s/exper/lxy/code/BEApp/*File.log
rm -r /home/k8s/exper/lxy/code/BEApp/resourceusage.log


# sleep 360

while :
do
mem=`kubectl top pod |grep mongo-0 |awk '{print $3}'| tr -cd "[0-9]"`

gi=1024

mem=$[ $mem / $gi ]

echo $mem  >> /home/k8s/exper/lxy/code/BEApp/resourceusage.log

for((j=1;j<=12;j++));
do

spark_master=`kubectl get pod |grep spark-hadoop-hibench-master | awk '{print $1}'`

index=0
addIndex=1
for element in ${spark_master[@]}
do
    temp_master[$index]=$element
    index=$[$index + $addIndex]
done


for (( i=0;i<${#temp_master[@]};i++)) do

    spark_master[$i]=${temp_master[$i]}

done


for (( i=0;i<${#spark_master[@]};i++)) do

spark_master[$i]=${spark_master[$i]: 28: 1}

done

echo ${spark_master[@]}

nowtime=`date +'%H:%M:%S'`


for element in ${spark_master[@]}
do

echo $element

echo -n "$nowtime " >> /home/k8s/exper/lxy/code/BEApp/spark-hadoop-hibench-master-$element-0File.log

kubectl exec spark-hadoop-hibench-slave-$element-0 -- du -k  /root/spark-hadoop/hadoop-2.7.3/tmp/nm-local-dir/usercache/root/appcache |awk 'END{print}'|awk '{print$1}'  >>/home/k8s/exper/lxy/code/BEApp/spark-hadoop-hibench-master-$element-0File.log


done

sleep 5

done

done