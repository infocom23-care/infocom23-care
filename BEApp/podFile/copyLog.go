package podFile

import "os/exec"

func CopyLog(podName string) {

	var instructions string = "kubectl cp " + podName + ":/root/log/spark.log /home/k8s/exper/lxy/code/BEApp/" + podName + ".log"
	//fmt.Printf("%s\n", instructions)

	cmd := exec.Command("sh", "-c", instructions)

	cmd.Run()
}
