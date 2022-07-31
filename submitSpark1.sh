./startCluster.sh $1 $2

sleep 30

./initHibench.sh $2

nowtime=`date +'%Y-%m-%d %H:%M:%S'`
echo "集群1实验开始时间： $nowtime" >> testAll.log

starttime=`date +'%Y-%m-%d %H:%M:%S'`
start_seconds=$(date --date="$starttime" +%s);

kubectl exec spark-hadoop-hibench-master-1-0 -- /root/spark-hadoop/HiBench-HiBench-7.1/bin/workloads/websearch/pagerank/spark/run.sh 

endtime=`date +'%Y-%m-%d %H:%M:%S'`
end_seconds=$(date --date="$endtime" +%s);
echo "集群1实验结束时间： $endtime" >> testAll.log
echo "集群1运行时间： "$((end_seconds-start_seconds))"s" >> testAll.log

sleep 120

kubectl delete svc spark-hadoop-hibench-master-$2

kubectl delete svc spark-hadoop-hibench-slave-$2

kubectl delete statefulsets.apps spark-hadoop-hibench-master-$2

kubectl delete statefulsets.apps spark-hadoop-hibench-slave-$2