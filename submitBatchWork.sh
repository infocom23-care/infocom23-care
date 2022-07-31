#!/bin/bash

cluster_nums=(1 2 3 4 5 6 7 8)
memory_arry=(3 3 2 3 3 3 2 3)
replicas=4
master=4

MEMGbyte=67

# MEMGbyte=77

while :
do

memGbyte=`ssh skv-node5 -- free -g | awk 'NR==2 {print $3}'`

remainMem=$[$MEMGbyte - $memGbyte]
echo "remainmem $remainMem"

# length=${#cluster_nums[@]}

# if [ $length = "0" ]; then

#     kubectl uncordon skv-node5

# fi

for (( i=0;i<${#memory_arry[@]};i++)) 
do

    total_mem=$[$replicas*${memory_arry[$i]}]
    total_mem=$[$total_mem+$master]

    echo "total mem $total_mem"

    if [[ "$remainMem" -ge "$total_mem" ]];
    then

        echo 'sucess'

        echo "choose ${cluster_nums[$i]}"

        ./submitSpark${cluster_nums[$i]}.sh $replicas ${cluster_nums[$i]} &

        sleep 60

        kubectl cordon skv-node5

        unset cluster_nums[$i]
        unset memory_arry[$i]

        index=0
        addIndex=1
        tempArry_mem=()
        tempArry_cluster=()

        for element in ${memory_arry[@]}
        do
        tempArry_mem[$index]=$element

        index=$[$index + $addIndex]
        done


        index=0
        for element in ${cluster_nums[@]}
        do
        tempArry_cluster[$index]=$element

        index=$[$index + $addIndex]
        done

        unset memory_arry
        unset cluster_nums

        echo ${tempArry_mem[0]}

        for (( i=0;i<${#tempArry_cluster[@]};i++)) do

            # echo ${tempArry_cluster[$i]}

            memory_arry[$i]=${tempArry_mem[$i]}
            cluster_nums[$i]=${tempArry_cluster[$i]}

        done

        echo ${cluster_nums[@]}
        echo ${memory_arry[@]}
        break

    fi

done

spark_master=`kubectl get pod |grep spark-hadoop-hibench-master | awk '{print $1}'`
master_style=`kubectl get pod |grep spark-hadoop-hibench-master | awk '{print $3}'`

# echo ${spark_master[@]}
# echo ${master_style[@]}

temp_master=()
temp_style=()
index=0
addIndex=1
for element in ${spark_master[@]}
do
    temp_master[$index]=$element
    index=$[$index + $addIndex]
done

index=0
for element in ${master_style[@]}
do
    temp_style[$index]=$element
    index=$[$index + $addIndex]
done

spark_master=()
master_style=()

for (( i=0;i<${#temp_master[@]};i++)) do

    spark_master[$i]=${temp_master[$i]}
    master_style[$i]=${temp_style[$i]}

done


for (( i=0;i<${#spark_master[@]};i++)) do

spark_master[$i]=${spark_master[$i]: 28: 1}

done

length=${#cluster_nums[@]}

for (( i=0;i<${#master_style[@]};i++)) do

#    echo ${master_style[$i]}
  if [ ${master_style[$i]} = "Pending" ] || [ ${master_style[$i]} = "Evicted" ];then
    echo "sucess"

    cluster_nums[$length]=${spark_master[$i]}

    if [ ${spark_master[$i]} = "3" ] || [ ${spark_master[$i]} = "7" ];then

        memory_arry[$length]=2

    else
        memory_arry[$length]=3
    fi

    length=$[$length + $addIndex]

    kubectl delete svc spark-hadoop-hibench-master-${spark_master[$i]}

    kubectl delete svc spark-hadoop-hibench-slave-${spark_master[$i]}

    kubectl delete statefulsets.apps spark-hadoop-hibench-master-${spark_master[$i]}

    kubectl delete statefulsets.apps spark-hadoop-hibench-slave-${spark_master[$i]}

  fi

done

echo ${spark_master[@]}
echo ${master_style[@]}

echo ${cluster_nums[@]}
echo ${memory_arry[@]}


# cluster_nums=3
# replicas=4

# ./startCluster.sh $replicas $cluster_nums

# sleep 10

# ./initHibench.sh $cluster_nums

done