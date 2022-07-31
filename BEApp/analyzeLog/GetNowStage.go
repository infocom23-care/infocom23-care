package analyzeLog

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetNowStage(podName string, SparkNowStage map[string]int, SparkNowData map[string]float64, SparkNowStageProgress map[string]float64) {
	filePath := "/home/k8s/exper/lxy/code/BEApp/" + podName + ".log"

	//fmt.Printf("%s\n", filePath)

	f, err := os.Open(filePath)

	var nowData float64
	nowData = 0.0

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {

		if strings.Contains(scanner.Text(), "stage") && strings.Contains(scanner.Text(), "Finished") {
			//fmt.Println(scanner.Text())
			tempVector := strings.Fields(scanner.Text())

			Length := len(tempVector)

			//fmt.Println(tempVector[Length-1])

			processTempVector := strings.Split(tempVector[Length-1], "/")

			nowProcess, _ := strconv.ParseFloat(processTempVector[0][1:], 64)
			totalProcess, _ := strconv.ParseFloat(processTempVector[1][0:len(processTempVector[1])-1], 64)

			SparkNowStageProgress[podName] = nowProcess / totalProcess

			//fmt.Println(SparkNowStageProgress[podName])
			//fmt.Println(nowProcess)
			//fmt.Println(totalProcess)
			//fmt.Println(processTempVector[0])

			for index := 0; index < Length; index++ {
				//fmt.Println(tempVector[index])
				if strings.EqualFold(tempVector[index], "stage") && index+1 <= Length {

					// fmt.Println(tempVector[index+1])

					resultString := strings.Split(tempVector[index+1], ".")
					//fmt.Println(resultString[0])

					SparkNowStage[podName], _ = strconv.Atoi(resultString[0])
				}
			}
		}

		if strings.Contains(scanner.Text(), "Finished") && strings.Contains(scanner.Text(), "stage") {
			//fmt.Println(scanner.Text())
			tempVector := strings.Fields(scanner.Text())

			//Length := len(tempVector)

			//fmt.Println(tempVector[13])

			tempData, _ := strconv.ParseFloat(tempVector[13], 64)
			tempData = tempData / 1000

			//fmt.Println(tempData)

			nowData += tempData

			continue
		}

	}

	SparkNowData[podName] = nowData
}

func GetStartTime(podName string) (int, int, int) {
	filePath := "/home/k8s/exper/lxy/code/BEApp/" + podName + ".log"

	//fmt.Printf("%s\n", filePath)

	startSecond := 0
	startMintu := 0
	startHour := 0

	f, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return startSecond, startMintu, startHour
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		tempVector := strings.Fields(scanner.Text())

		//fmt.Println(tempVector[1])

		if !strings.EqualFold(tempVector[2], "INFO") {
			continue
		}

		tempTimeVector := strings.Split(tempVector[1], ":")

		//fmt.Println(tempTimeVector[0])
		startSecond, _ = strconv.Atoi(tempTimeVector[2])
		startMintu, _ = strconv.Atoi(tempTimeVector[1])
		startHour, _ = strconv.Atoi(tempTimeVector[0])

		// fmt.Println(startMintu)
		// fmt.Println(startHour)

		break
	}

	return startSecond, startMintu, startHour

}

func GetEndTime(podName string) (int, int, int) {
	filePath := "/home/k8s/exper/lxy/code/BEApp/" + podName + ".log"

	//fmt.Printf("%s\n", filePath)

	endSecond := 0
	endMintu := 0
	endHour := 0

	f, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return endSecond, endMintu, endHour
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		tempVector := strings.Fields(scanner.Text())

		//fmt.Println(tempVector[1])

		if len(tempVector) < 3 {
			continue
		}

		if !strings.EqualFold(tempVector[2], "INFO") {
			continue
		}

		tempTimeVector := strings.Split(tempVector[1], ":")

		//fmt.Println(tempTimeVector[0])
		endSecond, _ = strconv.Atoi(tempTimeVector[2])
		endMintu, _ = strconv.Atoi(tempTimeVector[1])
		endHour, _ = strconv.Atoi(tempTimeVector[0])

		// fmt.Println(startMintu)
		// fmt.Println(startHour)

	}

	return endSecond, endMintu, endHour
}

func GetUnCompletedTask(podName string) (int, int) {
	filePath := "/home/k8s/exper/lxy/code/BEApp/" + podName + ".log"

	//fmt.Printf("%s\n", filePath)

	f, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return 0, 0
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	startTask := 0

	finishTask := 0

	for scanner.Scan() {

		if strings.Contains(scanner.Text(), "stage") && strings.Contains(scanner.Text(), "Finished") {
			//fmt.Println(scanner.Text())
			finishTask += 1
		}

		if strings.Contains(scanner.Text(), "Starting") && strings.Contains(scanner.Text(), "stage") {
			startTask += 1
		}

	}
	return startTask, startTask - finishTask
}

func GetRunTimeUnCompleted(SparkNowData map[string]float64, startTask int, unCompleteTask int, podName string) {

	startIndex := startTask - unCompleteTask + 1

	nowIndex := 0

	filePath := "/home/k8s/exper/lxy/code/BEApp/" + podName + ".log"

	//fmt.Printf("%s\n", filePath)

	f, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	startSecond := 0
	startMintu := 0
	startHour := 0

	for scanner.Scan() {

		if strings.Contains(scanner.Text(), "Starting") && strings.Contains(scanner.Text(), "stage") {
			nowIndex += 1
		}

		if nowIndex == startIndex {

			tempVector := strings.Fields(scanner.Text())

			//fmt.Println(tempVector[1])

			tempTimeVector := strings.Split(tempVector[1], ":")

			//fmt.Println(tempTimeVector[0])
			startSecond, _ = strconv.Atoi(tempTimeVector[2])
			startMintu, _ = strconv.Atoi(tempTimeVector[1])
			startHour, _ = strconv.Atoi(tempTimeVector[0])

			hours, minutes, second := time.Now().Clock()

			nowSpendTime := (hours-startHour)*60*60 - (startMintu-minutes)*60 - (startSecond - second)

			SparkNowData[podName] += float64(nowSpendTime)

			startIndex++

			if startIndex > startIndex+unCompleteTask {
				break
			}

		}

	}

}

func GetTotalCpuCost(podName string, SparkNowData map[string]float64) {
	startTask, unCompleteTask := GetUnCompletedTask(podName)
	GetRunTimeUnCompleted(SparkNowData, startTask, unCompleteTask, podName)

}

func GetShuffle(podName string, SparkShuffleData map[string]float64) {
	filePath := "/home/k8s/exper/lxy/code/BEApp/" + podName + ".log"

	//fmt.Printf("%s\n", filePath)

	f, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	Second := 0
	Mintu := 0
	Hour := 0

	flag := 0

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Asked") {

			flag = 1
			continue
		}

		if flag == 1 {
			tempVector := strings.Fields(scanner.Text())

			//fmt.Println(tempVector[1])

			tempTimeVector := strings.Split(tempVector[1], ":")

			//fmt.Println(tempTimeVector[0])
			Second, _ = strconv.Atoi(tempTimeVector[2])
			Mintu, _ = strconv.Atoi(tempTimeVector[1])
			Hour, _ = strconv.Atoi(tempTimeVector[0])

			flag = 0
		}
	}

	// fmt.Print("时间" + podName + " ")
	// fmt.Print(Hour)
	// fmt.Print("  ")
	// fmt.Print(Mintu)
	// fmt.Print("  ")
	// fmt.Print(Second)
	// fmt.Print("  ")
	// fmt.Println()

	Second = Second - Second%10

	// if Second == 0 {

	// 	Second = 50

	// 	Mintu = (Mintu - 1) % 60

	// 	if Mintu == 59 {
	// 		Hour = Hour - 1
	// 	}
	// } else {

	// 	Second = Second - 10

	// }

	// fmt.Print("时间" + podName + " ")
	// fmt.Print(Hour)
	// fmt.Print("  ")
	// fmt.Print(Mintu)
	// fmt.Print("  ")
	// fmt.Print(Second)
	// fmt.Print("  ")
	// fmt.Println()

	targetSecond := 0
	targetMintu := 0
	targetHour := 0

	filePath = "/home/k8s/exper/lxy/code/BEApp/" + podName + "File.log"

	// fmt.Printf("%s\n", filePath)

	f, err = os.Open(filePath)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer f.Close()

	scanner = bufio.NewScanner(f)

	for scanner.Scan() {
		tempVector := strings.Fields(scanner.Text())

		//fmt.Println(tempVector[1])

		tempTimeVector := strings.Split(tempVector[0], ":")

		//fmt.Println(tempTimeVector[0])
		targetSecond, _ = strconv.Atoi(tempTimeVector[2])
		targetMintu, _ = strconv.Atoi(tempTimeVector[1])
		targetHour, _ = strconv.Atoi(tempTimeVector[0])

		targetSecond = targetSecond - targetSecond%10

		if Second == targetSecond && Mintu == targetMintu && Hour == targetHour {

			tempNum, _ := strconv.Atoi(tempVector[1])

			SparkShuffleData[podName] = float64(tempNum) / float64(1024)

		}
	}

}

func GetProgressForecast(podName string, SparkNowStage map[string]int, SparkAllStage map[string]int, SparkNowStageProgress map[string]float64, SparkNowData map[string]float64, SparkRemainTime map[string]float64) {

	nowSpendTime := 0

	startSecond, startMintu, startHour := GetStartTime(podName)
	endSecond, endMintu, endHour := GetEndTime(podName)

	nowSpendTime = (endHour-startHour)*60*60 + (-startMintu+endMintu)*60 + (-startSecond + endSecond)

	// fmt.Print(podName)
	// fmt.Print("当前花费：")
	// fmt.Println(nowSpendTime)

	processOfStage := (float64((SparkNowStage[podName])) + SparkNowStageProgress[podName]) / float64((SparkAllStage[podName] + 1.0))

	remainProcess := 1 - processOfStage
	fmt.Print(podName)
	fmt.Print("剩余阶段")
	fmt.Println(remainProcess)

	remainTime := 3.0 * float64(nowSpendTime) / float64(processOfStage) * (remainProcess)

	fmt.Print("剩余时间:")
	fmt.Println(remainTime)

	SparkRemainTime[podName] = remainTime
}
