package analyzeLog

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func GetVariance(w int) float64 {

	var ResourceUsage map[int]float64 = make(map[int]float64)

	filePath := "/home/k8s/exper/lxy/code/BEApp/resourceusage.log"
	f, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return 0.0
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var i int = 1

	for scanner.Scan() {

		tempResource, _ := strconv.ParseFloat(scanner.Text(), 64)
		ResourceUsage[i%w] = tempResource
		i++

	}

	var sum float64 = 0.0
	var len int = 0
	var result float64 = 0.0

	for _, v := range ResourceUsage {

		sum += v
		len++
	}

	avager := sum / float64(len)

	for _, v := range ResourceUsage {

		result += (avager - v) * (avager - v)
	}

	// for k, v := range ResourceUsage {
	// 	fmt.Printf("resource useage")
	// 	fmt.Printf("%d   %f\n", k, v)
	// }

	result = result / float64(len)

	// fmt.Print("方差是")
	// fmt.Println(result)
	// fmt.Println(len)

	return result

}
