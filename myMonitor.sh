#!/bin/bash
echo $podNmae

MEMGbyte=68

#MEMGbyte=78

while :
do

memGbyte=`ssh 10.0.0.65 -- free -g | awk 'NR==2 {print $3}'`
echo $memGbyte

if [[ "$memGbyte" -ge "$MEMGbyte" ]]; then
    /home/k8s/exper/lxy/code/BEApp/./main 

    chmod u+x /home/k8s/exper/lxy/code/BEApp/result.txt

fi

while [[ "$memGbyte" -ge "$MEMGbyte" ]]; do

    podname=`cat /home/k8s/exper/lxy/code/BEApp/result.txt |head -n 1`

    sed -i '1d' /home/k8s/exper/lxy/code/BEApp/result.txt

    echo $podname >> testAll.log
    kubectl delete pod $podname
    sleep 5
    memGbyte=`ssh skv-node5 -- free -g | awk 'NR==2 {print $3}'`
    echo $memGbyte
done

done

echo "all done"