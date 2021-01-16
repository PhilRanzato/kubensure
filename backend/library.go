package backend

import (
	"log"
	"os"
	"path/filepath"

	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

// GetPodByName : accepts a list of pods and a pod name+namespace
//			returns a Pod
func GetPodByName(pods []v1.Pod, podName string, podNamespace string) v1.Pod {
	var pod v1.Pod
	for _, p := range pods {
		if p.Name == podName && p.Namespace == podNamespace {
			pod = p
		}
	}
	return pod
}

// GetPodIP : returns the pod ip
func GetPodIP(pod v1.Pod) string {
	return pod.Status.PodIP
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

// GetServiceByName : accepts a list of services and a service name+namespace
//			returns a Service
func GetServiceByName(svcs []v1.Service, svcName string, svcNamespace string) v1.Service {
	var svc v1.Service
	for _, s := range svcs {
		if s.Name == svcName && s.Namespace == svcNamespace {
			svc = s
		}
	}
	return svc
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
