package main

import (
	"BEApp/analyzeLog"
	"math"

	// "BEApp/podFile"
	"BEApp/podFile"
	"context"
	"flag"
	"fmt"
	"os/exec"

	// "os/exec"

	// "os/exec"
	"path/filepath"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	start := time.Now()

	var SparkPodMap map[string]int = make(map[string]int)

	var SparkAllStage map[string]int = make(map[string]int)
	var SparkNowStage map[string]int = make(map[string]int)
	var SparkNowData map[string]float64 = make(map[string]float64)
	var SparkShuffleData map[string]float64 = make(map[string]float64)
	var SparkNowStageProgress map[string]float64 = make(map[string]float64)
	var SparkRemainTime map[string]float64 = make(map[string]float64)

	var SparAllCost map[string]float64 = make(map[string]float64)

	var kubeconfig *string

	// 读取配置文件
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err)
	}

	pod, _ := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

	for _, tempPod := range pod.Items {

		if strings.Contains(tempPod.Name, "spark-hadoop-hibench-master") {
			SparkPodMap[tempPod.Name] = 0
			SparkAllStage[tempPod.Name] = -1
			SparkNowStage[tempPod.Name] = -1
			SparkNowData[tempPod.Name] = 0.0
			SparkPodMap[tempPod.Name] = 0.0
			SparkNowStageProgress[tempPod.Name] = 0.0
			SparkShuffleData[tempPod.Name] = 0.0
			SparAllCost[tempPod.Name] = math.SmallestNonzeroFloat64
			SparkRemainTime[tempPod.Name] = 0.0
		}

		for key, _ := range SparkPodMap {
			//fmt.Printf(string(key[len(key)-3 : len(key)-2]))
			if string(key[len(key)-3:len(key)-2]) == string(tempPod.Name[len(tempPod.Name)-3:len(tempPod.Name)-2]) && strings.Contains(tempPod.Name, "spark-hadoop-hibench-slave") {
				//fmt.Printf("%s\n", key)
				SparkPodMap[key]++
			}
		}
	}

	exec.Command("sh", "-c", "rm -r /home/k8s/exper/lxy/code/BEApp/spark-hadoop-hibench-master-*-0.log").Run()

	for k, _ := range SparkPodMap {
		podFile.CopyLog(k)
	}

	for k, _ := range SparkPodMap {
		//fmt.Println(analyzeLog.GetResultStage(k))
		analyzeLog.GetResultStage(k, SparkAllStage)
	}

	for k, _ := range SparkPodMap {
		//fmt.Println(analyzeLog.GetResultStage(k))
		analyzeLog.GetNowStage(k, SparkNowStage, SparkNowData, SparkNowStageProgress)
	}

	for k, _ := range SparkPodMap {
		//fmt.Println(analyzeLog.GetResultStage(k))
		analyzeLog.GetTotalCpuCost(k, SparkNowData)
	}

	for k, _ := range SparkPodMap {
		//fmt.Println(analyzeLog.GetResultStage(k))
		analyzeLog.GetShuffle(k, SparkShuffleData)
	}

	for k, _ := range SparkPodMap {
		//fmt.Println(analyzeLog.GetResultStage(k))
		analyzeLog.GetProgressForecast(k, SparkNowStage, SparkAllStage, SparkNowStageProgress, SparkNowData, SparkRemainTime)
	}

	var w int = 4

	analyzeLog.GetAllCost(SparkShuffleData, SparkNowData, SparAllCost, SparkPodMap, analyzeLog.GetVariance(w), SparkRemainTime)

	// for k, v := range SparkAllStage {
	// 	fmt.Printf("%s   %d\n", k, v)
	// }

	// for k, v := range SparkNowStage {
	// 	fmt.Printf("%s   %d\n", k, v)
	// }

	for k, v := range SparkNowData {
		fmt.Printf("cpu cost")
		fmt.Printf("%s   %f\n", k, v)
	}

	// for k, v := range SparkNowStageProgress {
	// 	fmt.Printf("%s   %f\n", k, v)
	// }

	for k, v := range SparkShuffleData {
		fmt.Printf("shuffle data")
		fmt.Printf("%s   %f\n", k, v)
	}

	cost := time.Since(start)
	fmt.Printf("cost=[%s]\n", cost)
}
