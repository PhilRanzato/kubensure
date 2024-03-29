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

type networkCommandStructure struct {
	command       string
	portMandatory bool
}

type networkCommand struct {
	command string
}

var networkCommandConstructorList = []networkCommandStructure{
	networkCommandStructure{"wget --spider -q --timeout=5 %s", false},
	networkCommandStructure{"curl -s -k %s", false},
	networkCommandStructure{"nmap -p %s %d", true},
	networkCommandStructure{"nc -z -v -w 2 %s %d", true},
	networkCommandStructure{"echo -n | telnet %s %d", true},
}

func networkCommandConstructor(ncs networkCommandStructure, ep string, epNs string, port int) networkCommand {

	var nc networkCommand

	if epNs != "" {
		epNs = "." + epNs
	}

	if ncs.portMandatory && port != 0 {
		nc = networkCommand{
			command: fmt.Sprintf(ncs.command, ep+epNs, port),
		}
	} else if ncs.portMandatory && port == 0 {
	} else if port != 0 {
		nc = networkCommand{
			command: fmt.Sprintf(ncs.command+":%d", ep+epNs, port),
		}
	} else {
		nc = networkCommand{
			command: fmt.Sprintf(ncs.command, ep+epNs),
		}
	}

	return nc
}

func networkCommandList(ep string, epNs string, port int) []networkCommand {
	var ncl []networkCommand

	for _, t := range networkCommandConstructorList {
		ncl = append(ncl, networkCommandConstructor(t, ep, epNs, port))
	}
	return ncl
}

// ExecIntoPod : accepts a clientset, a pod, a command and a standard redader
//				 executes the specified command into the specified pod
func ExecIntoPod(clientset *kubernetes.Clientset, pod *v1.Pod, command string, stdin io.Reader, getOnlyResultCode bool) (string, string, error) {
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
	if getOnlyResultCode == false {
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

// ConnectionPodToService : accepts a pod and a service
//				 executes the specified command into the specified pod to test connection to the specified service
func ConnectionPodToService(clientset *kubernetes.Clientset, pod v1.Pod, svc v1.Service, svcPort int) bool {

	var result = false
	commands := networkCommandList(svc.Name, svc.Namespace, svcPort)
	for _, nc := range commands {
		if nc.command != "" {
			fmt.Printf("Testing with %s\n", nc.command)
			_, stderr, err := ExecIntoPod(clientset, &pod, nc.command, nil, true)
			if len(stderr) != 0 {
				// fmt.Println("STDERR:", stderr)
			}
			if err != nil {
				// fmt.Printf("Error occured while `exec`ing to the Pod %s, namespace %s, command %s\n", pod.Name, pod.Namespace, nc.command)
				// fmt.Println(err)
			} else {
				// fmt.Println("Output:")
				// fmt.Println(output)
				result = true
				break
			}
		}
	}
	return result
}

// ConnectionPodToExternal : accepts a pod and an external endpoint
//				 executes the specified command into the specified pod to test connection to the specified external endpoint
func ConnectionPodToExternal(clientset *kubernetes.Clientset, pod v1.Pod, url string, urlPort int) bool {

	var result = false
	commands := networkCommandList(url, "", urlPort)
	for _, nc := range commands {
		if nc.command != "" {
			fmt.Printf("Testing with %s\n", nc.command)
			_, stderr, err := ExecIntoPod(clientset, &pod, nc.command, nil, true)
			if len(stderr) != 0 {
				// fmt.Println("STDERR:", stderr)
			}
			if err != nil {
				// fmt.Printf("Error occured while `exec`ing to the Pod %s, namespace %s, command %s\n", pod.Name, pod.Namespace, nc.command)
				// fmt.Println(err)
			} else {
				// fmt.Println("Output:")
				// fmt.Println(output)
				result = true
				break
			}
		}
	}
	return result
}

// ConnectionPodToPod : accepts two pods
//				 executes the specified command into the specified pod to test connection against the other pod
func ConnectionPodToPod(clientset *kubernetes.Clientset, pod v1.Pod, target v1.Pod, targetPort int) bool {

	var result = false
	commands := networkCommandList(strings.ReplaceAll(GetPodIP(target), ".", "-"), target.Namespace+".pod", targetPort)
	for _, nc := range commands {
		if nc.command != "" {
			fmt.Printf("Testing with %s\n", nc.command)
			_, stderr, err := ExecIntoPod(clientset, &pod, nc.command, nil, true)
			if len(stderr) != 0 {
				// fmt.Println("STDERR:", stderr)
			}
			if err != nil {
				// fmt.Printf("Error occured while `exec`ing to the Pod %s, namespace %s, command %s\n", pod.Name, pod.Namespace, nc.command)
				// fmt.Println(err)
			} else {
				// fmt.Println("Output:")
				// fmt.Println(output)
				result = true
				break
			}
		}
	}
	return result
}
