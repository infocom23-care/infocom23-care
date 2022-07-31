package analyzeLog

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

func GetAllCost(SparkShuffleData map[string]float64, SparkNowData map[string]float64, SparkAllCost map[string]float64, SparkPodMap map[string]int, variance float64, SparkRemainTime map[string]float64) {

	var instructions string = "rm /home/k8s/exper/lxy/code/BEApp/result.txt"
	//fmt.Printf("%s\n", instructions)

	cmd := exec.Command("sh", "-c", instructions)

	cmd.Run()

	// fmt.Println("success")

	a := 1.0
	b := 0.05
	c := 0.005

	// fmt.Print("方差是：")
	// fmt.Println(variance)

	for k, _ := range SparkNowData {

		SparkAllCost[k] = a*SparkNowData[k] + b*SparkShuffleData[k] - c*variance*SparkRemainTime[k]
	}

	type SparkToCost struct {
		sparkName string

		cost float64
	}

	var AllSparkCost []SparkToCost

	for k, v := range SparkAllCost {

		AllSparkCost = append(AllSparkCost, SparkToCost{k, v})
	}

	sort.Slice(AllSparkCost, func(i, j int) bool {
		return AllSparkCost[i].cost < AllSparkCost[j].cost
	})

	fmt.Println(AllSparkCost)

	fp := "/home/k8s/exper/lxy/code/BEApp/result.txt"

	f, err := os.OpenFile(fp, os.O_APPEND, 0660)
	defer f.Close()
	if err != nil {
		// fmt.Println("文件不存在,创建文件")
		f, _ = os.Create(fp)
	}

	for k := range AllSparkCost {

		slaveNum := SparkPodMap[AllSparkCost[k].sparkName]

		for j := 0; j < slaveNum; j++ {

			tempString := AllSparkCost[k].sparkName[0:21] + "slave-" + AllSparkCost[k].sparkName[28:29] + "-" + strconv.Itoa(SparkPodMap[AllSparkCost[k].sparkName]-j-1)

			// fmt.Println(tempString)

			f.WriteString(tempString)
			f.WriteString("\n")

		}

		f.WriteString(AllSparkCost[k].sparkName)
		f.WriteString("\n")

	}

	for k := range AllSparkCost {
		f.WriteString(AllSparkCost[k].sparkName)
		f.WriteString("  ")
		tempString := SparkPodMap[AllSparkCost[k].sparkName]
		f.WriteString(strconv.Itoa(tempString))
		f.WriteString("\n")
	}

}
