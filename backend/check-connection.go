package backend

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
)

type networkExecutable struct {
	command, args string
	needPort      bool
}

var networkExecutables = []networkExecutable{
	networkExecutable{
		command:  "wget",
		args:     "--spider --timeout=1",
		needPort: false,
	},
	networkExecutable{
		command:  "nslookup",
		args:     "",
		needPort: false,
	},
	networkExecutable{
		command:  "curl",
		args:     "-O",
		needPort: false,
	},
	networkExecutable{
		command:  "nc",
		args:     "",
		needPort: true,
	},
	networkExecutable{
		command:  "telnet",
		args:     "",
		needPort: true,
	},
	networkExecutable{
		command:  "nmap",
		args:     "-p",
		needPort: true,
	},
	networkExecutable{
		command:  "ping",
		args:     "",
		needPort: false,
	},
	networkExecutable{
		command: "dig",
		args:    "",
	},
}

// ExecIntoPod : accepts a clientset, a pod, a command and a standard redader
//				 executes the specified command into the specified pod
func ExecIntoPod(clientset *kubernetes.Clientset, pod *v1.Pod, command string, stdin io.Reader) (string, string, error) {
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := v1.AddToScheme(scheme); err != nil {
		return "", "", fmt.Errorf("error adding to scheme: %v", err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&v1.PodExecOptions{
		Command: strings.Fields(command),
		Stdin:   stdin != nil,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(GetConfig(), "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("error while creating Executor: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return "", "", fmt.Errorf("error in Stream: %v", err)
	}

	return stdout.String(), stderr.String(), nil

}

func getNetworkExecutable(pod v1.Pod) (networkExecutable, bool) {

	clientset := GetClientSet()
	var command string
	var result networkExecutable

	for idx, ne := range networkExecutables {
		if ne.needPort {
			command = ne.command
		} else {
			command = ne.command
		}

		output, stderr, err := ExecIntoPod(clientset, &pod, command, nil)

		if len(stderr) != 0 {
			fmt.Println("STDERR:", stderr)
		}
		if err != nil {
			fmt.Printf("Error occured while `exec`ing to the Pod %s, namespace %s, command %s (executable %s)\n", pod.Name, pod.Namespace, command, networkExecutables[idx].command)
			fmt.Println(err)
		} else {
			fmt.Println("Output:")
			fmt.Println(output)
			result = ne
		}
	}

	if result == (networkExecutable{}) {
		return result, false
	}

	return result, true
}

// TestConnectionPodToService : accepts a pod and a service
//				 executes the specified command into the specified pod to test connection to the specified service
func TestConnectionPodToService(pod v1.Pod, svc v1.Service) bool {

	clientset := GetClientSet()

	// execute command
	var result bool
	command := "wget --spider --timeout=1 " + svc.Name + "." + svc.Namespace
	output, stderr, err := ExecIntoPod(clientset, &pod, command, nil)
	if len(stderr) != 0 {
		fmt.Println("STDERR:", stderr)
		result = false
	}
	if err != nil {
		fmt.Printf("Error occured while `exec`ing to the Pod %s, namespace %s, command %s\n", pod.Name, pod.Namespace, command)
		fmt.Println(err)
		result = false
	} else {
		fmt.Println("Output:")
		fmt.Println(output)
		result = true
	}

	return result
}
