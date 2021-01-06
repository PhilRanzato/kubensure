package backend

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type networkExecutable struct {
	command, args string
}

var networkExecutables = []networkExecutable{
	networkExecutable{
		command: "wget",
		args:    "--spider --timeout=1",
	},
	networkExecutable{
		command: "nslookup",
		args:    "",
	},
	networkExecutable{
		command: "curl",
		args:    "-O",
	},
	networkExecutable{
		command: "nc",
		args:    "-O",
	},
	networkExecutable{
		command: "telnet",
		args:    "",
	},
	networkExecutable{
		command: "nmap",
		args:    "-p",
	},
	networkExecutable{
		command: "ping",
		args:    "",
	},
	networkExecutable{
		command: "dig",
		args:    "",
	},
}

// Getters

// GetConfig : Get kubeconfig from file
func GetConfig() *rest.Config {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

// GetClientSet : get client set from the kubeconfig file
func GetClientSet() *kubernetes.Clientset {
	var config = GetConfig()

	// Create a rest client not targeting specific API version
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}

// GetPods : accepts a clientset and returns a list of Pods
func GetPods(clientset *kubernetes.Clientset) []v1.Pod {
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get Pods:", err)
	}
	return pods.Items
}

// GetDeployments : accepts a clientset and returns a list of Deployments
func GetDeployments(clientset *kubernetes.Clientset) []appv1.Deployment {
	deploys, err := clientset.AppsV1().Deployments("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get Deployments:", err)
	}
	return deploys.Items
}

// GetDaemonSets : accepts a clientset and returns a list of DaemonSets
func GetDaemonSets(clientset *kubernetes.Clientset) []appv1.DaemonSet {
	ds, err := clientset.AppsV1().DaemonSets("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get DaemonSets:", err)
	}
	return ds.Items
}

// GetReplicaSets : accepts a clientset and returns a list of ReplicaSets
func GetReplicaSets(clientset *kubernetes.Clientset) []appv1.ReplicaSet {
	rs, err := clientset.AppsV1().ReplicaSets("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get ReplicaSets:", err)
	}
	return rs.Items
}

// GetStatefulSets : accepts a clientset and returns a list of StatefulSets
func GetStatefulSets(clientset *kubernetes.Clientset) []appv1.StatefulSet {
	sts, err := clientset.AppsV1().StatefulSets("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get StatefulSets:", err)
	}
	return sts.Items
}

// GetServices : accepts a clientset and returns a list of Services
func GetServices(clientset *kubernetes.Clientset) []v1.Service {
	svcs, err := clientset.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get Services:", err)
	}
	return svcs.Items
}

// GetSecrets : accepts a clientset and returns a list of Secrets
func GetSecrets(clientset *kubernetes.Clientset) []v1.Secret {
	scr, err := clientset.CoreV1().Secrets("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get Secrets:", err)
	}
	return scr.Items
}

// GetConfigMaps : accepts a clientset and returns a list of ConfigMaps
func GetConfigMaps(clientset *kubernetes.Clientset) []v1.ConfigMap {
	cm, err := clientset.CoreV1().ConfigMaps("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get ConfigMaps:", err)
	}
	return cm.Items
}

// GetServiceAccounts : accepts a clientset and returns a list of ServiceAccounts
func GetServiceAccounts(clientset *kubernetes.Clientset) []v1.ServiceAccount {
	sa, err := clientset.CoreV1().ServiceAccounts("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get ServiceAccounts:", err)
	}
	return sa.Items
}

// GetEvents : accepts a clientset and returns a list of Events
func GetEvents(clientset *kubernetes.Clientset) []v1.Event {
	ev, err := clientset.CoreV1().Events("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get Events:", err)
	}
	return ev.Items
}

// GetEndpoints : accepts a clientset and returns a list of Endpoints
func GetEndpoints(clientset *kubernetes.Clientset) []v1.Endpoints {
	ep, err := clientset.CoreV1().Endpoints("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get Endpoints:", err)
	}
	return ep.Items
}

// GetPersistentVolumes : accepts a clientset and returns a list of PersistentVolumes
func GetPersistentVolumes(clientset *kubernetes.Clientset) []v1.PersistentVolume {
	pvs, err := clientset.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get PersistentVolumes:", err)
	}
	return pvs.Items
}

// GetPersistentVolumeClaims : accepts a clientset and returns a list of PersistentVolumeClaims
func GetPersistentVolumeClaims(clientset *kubernetes.Clientset) []v1.PersistentVolumeClaim {
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get PersistentVolumeClaims:", err)
	}
	return pvcs.Items
}

// GetRoles : accepts a clientset and returns a list of Roles
func GetRoles(clientset *kubernetes.Clientset) []rbacv1.Role {
	roles, err := clientset.RbacV1().Roles("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get Roles:", err)
	}
	return roles.Items
}

// GetRoleBindings : accepts a clientset and returns a list of RoleBindings
func GetRoleBindings(clientset *kubernetes.Clientset) []rbacv1.RoleBinding {
	rb, err := clientset.RbacV1().RoleBindings("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get RoleBindings:", err)
	}
	return rb.Items
}

// GetClusterRoles : accepts a clientset and returns a list of ClusterRoles
func GetClusterRoles(clientset *kubernetes.Clientset) []rbacv1.ClusterRole {
	cr, err := clientset.RbacV1().ClusterRoles().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get ClusterRoles:", err)
	}
	return cr.Items
}

// GetClusterRoleBindings : accepts a clientset and returns a list of ClusterRoleBindings
func GetClusterRoleBindings(clientset *kubernetes.Clientset) []rbacv1.ClusterRoleBinding {
	crb, err := clientset.RbacV1().ClusterRoleBindings().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to get ClusterRoleBindings:", err)
	}
	return crb.Items
}

// GetServerVersion : accepts a clientset and returns a list of ServerVersion
func GetServerVersion(clientset *kubernetes.Clientset) string {
	version, err := clientset.ServerVersion()
	if err != nil {
		log.Fatalln("Failed to get ServerVersion:", err)
	}
	return version.String()
}

// Exec

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

// func getNetworkExecutable(pod v1.Pod) NetworkExecutable{

// }

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
