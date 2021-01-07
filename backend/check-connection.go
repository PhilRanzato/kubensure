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
		args:     " --spider --timeout=1",
		needPort: false,
	},
	networkExecutable{
		command:  "nslookup",
		args:     "",
		needPort: false,
	},
	networkExecutable{
		command:  "curl",
		args:     " -O",
		needPort: false,
	},
	networkExecutable{
		command:  "netcat",
		args:     "-q 0",
		needPort: true,
	},
	networkExecutable{
		command:  "nmap",
		args:     "-p",
		needPort: true,
	},
	// networkExecutable{
	// 	command:  "dig",
	// 	args:     "",
	// 	needPort: false,
	// },
}

// ExecIntoPod : accepts a clientset, a pod, a command and a standard redader
//				 executes the specified command into the specified pod
func ExecIntoPod(clientset *kubernetes.Clientset, pod *v1.Pod, command string, stdin io.Reader, testExecutable bool) (string, string, error) {
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
	if testExecutable == true {
		req.VersionedParams(&v1.PodExecOptions{
			Command: strings.Fields(command),
			Stdin:   stdin != nil,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
		}, parameterCodec)
	} else {
		req.VersionedParams(&v1.PodExecOptions{
			Command: []string{
				"sh",
				"-c",
				command + " &> /dev/null && echo $?",
			},
			Stdin:  stdin != nil,
			Stdout: true,
			Stderr: true,
			TTY:    false,
		}, parameterCodec)

	}

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

	for _, ne := range networkExecutables {

		command = ne.command
		_, stderr, err := ExecIntoPod(clientset, &pod, command, nil, true)

		if len(stderr) != 0 {
			fmt.Println("STDERR:", stderr)
		}
		if err != nil && err.Error() != "error in Stream: command terminated with exit code 1" {
			fmt.Printf("Executable %s not find in path (%s)\n", ne.command, err)
		} else {
			result = ne
		}
	}

	if result == (networkExecutable{}) {
		return result, false
	}

	return result, true
}

func buildNetworkCommand(ne networkExecutable, svc v1.Service, svcPort string) string {

	command := ""
	args := ""

	if len(ne.args) > 0 {
		args = " " + ne.args
	}
	if svc.Name != "" {
		if ne.needPort && svcPort != "" {
			command = ne.command + args + " " + svc.Name + "." + svc.Namespace + " " + svcPort
			// command = ne.command + args + " " + svc.Name + "." + svc.Namespace + " " + svcPort
		} else {
			command = ne.command + args + " " + svc.Name + "." + svc.Namespace
			// command = ne.command + args + " " + svc.Name + "." + svc.Namespace
		}
	} else {
		command = "exit 1"
	}

	return command
}

// TestConnectionPodToService : accepts a pod and a service
//				 executes the specified command into the specified pod to test connection to the specified service
func TestConnectionPodToService(clientset *kubernetes.Clientset, pod v1.Pod, svc v1.Service, svcPort string) bool {

	var result bool
	executable, present := getNetworkExecutable(pod)
	if present {
		command := buildNetworkCommand(executable, svc, svcPort)
		fmt.Printf("Testing with %s\n", command)
		output, stderr, err := ExecIntoPod(clientset, &pod, command, nil, false)
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
	} else {
		result = false
	}

	return result
}
