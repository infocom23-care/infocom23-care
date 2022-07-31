package analyzeLog

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetResultStage(podName string, SparkAllStage map[string]int) {

	filePath := "/home/k8s/exper/lxy/code/BEApp/" + podName + ".log"

	//fmt.Printf("%s\n", filePath)

	f, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		if strings.Contains(scanner.Text(), "ResultStage") {
			//fmt.Println(scanner.Text())
			tempVector := strings.Fields(scanner.Text())

			for index := 0; index < len(tempVector); index++ {
				//fmt.Println(tempVector[index])
				if strings.EqualFold(tempVector[index], "ResultStage") {
					result, _ := strconv.Atoi(tempVector[index+1])

					SparkAllStage[podName] = result
				}
			}

			break
		}

	}

}
