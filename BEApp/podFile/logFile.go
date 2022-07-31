package podFile

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func LogFile(config *rest.Config, clientset *kubernetes.Clientset, podName string) {

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace("default").
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   []string{"mkdir", "demotest"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
			Container: "spark-hibench-master",
		}, scheme.ParameterCodec)
	exec, _ := remotecommand.NewSPDYExecutor(config, "POST", req.URL())

	if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
		fmt.Errorf("stdin/stdout should be terminal")
	}

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		fmt.Println(err)
	}

	defer terminal.Restore(0, oldState)

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	if err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  screen,
		Stdout: screen,
		Stderr: screen,
		Tty:    false,
	}); err != nil {
		fmt.Print(err)
	}

}
