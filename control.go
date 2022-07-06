package main

import (
	//clientset "control-d/pkg/client/clientset/versioned"
	//informers "control-d/pkg/client/informers/externalversions"
	//v1 "control-d/pkg/client/informers/externalversions/openfaas/v1"
	"control-d/pkg/config"
	//"control-d/pkg/k8s"
	"flag"
	"fmt"
	"net/http"

	//"control-data/k8s"
	"context"
	//"control-d/pkg/signals"
	"io/ioutil"
	"log"

	//"math/rand"
	//"net/url"
	"os"
	"path/filepath"

	//"regexp"
	"strings"

	//"sync"
	//"time"

	providertypes "github.com/openfaas/faas-provider/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//kubeinformers "k8s.io/client-go/informers"
	//v1apps "k8s.io/client-go/informers/apps/v1"
	//v1core "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	//corelister "k8s.io/client-go/listers/core/v1"
	//"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

//var functionMatcher = regexp.MustCompile("^/?(?:async-)?function/([^/?]+)([^?]*)")

const (
	// hasPathCount = 3
	// routeIndex   = 0 // routeIndex corresponds to /function/ or /async-function/
	// nameIndex    = 1 // nameIndex is the function name
	// pathIndex    = 2 // pathIndex is the path i.e. /employee/:id/
	watchdogPort = 8080
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func control(functionName string) (resp *http.Response, err error) {

	//functionaddress := "/function/test-4"

	var apiGateway = "http://172.16.252.163:31112/function/"
	//var endpointsilices []string
	var servicesilices []string

	//functionName := getServiceName(functionaddress)

	fmt.Printf("function name is: %s \n", functionName)

	readConfig := config.ReadConfig{}
	osEnv := providertypes.OsEnv{}
	config, err := readConfig.Read(osEnv)

	namespace := "openfaas-fn"

	config.DefaultFunctionNamespace = namespace

	if err != nil {
		log.Fatalf("Error reading config: %s", err.Error())
	}

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	clientCmdConfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(clientCmdConfig)
	if err != nil {
		log.Fatalf("Error building Kubernetes clientset: %s", err.Error())
	}

	services, err := kubeClient.CoreV1().Services(config.DefaultFunctionNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, s := range services.Items {
		if strings.Contains(s.Name, functionName) {
			fmt.Printf("Name: %v Cluster IP: %v\n", s.Name, s.Spec.ClusterIP)
			servicesilices = append(servicesilices, s.Spec.ClusterIP)
		}

	}
	// pods, err := kubeClient.CoreV1().Pods(config.DefaultFunctionNamespace).List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// for _, v := range pods.Items {
	// 	if strings.Contains(v.Name, functionName) {
	// 		fmt.Printf("Name: %v IP: %v\n", v.Name, v.Status.PodIP)
	// 		endpointsilices = append(endpointsilices, v.Status.PodIP)
	// 	}

	// }
	if len(servicesilices) == 0 {
		resp, err := http.Get(apiGateway + functionName)
		defer resp.Body.Close()
		if err != nil {
			fmt.Printf("err")
		}
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("response is :%s", string(body))
		return resp, err
	} else {
		for i := range servicesilices {
			endpointIP := servicesilices[i]
			urlStr := fmt.Sprintf("http://%s:%d", endpointIP, watchdogPort)
			resp, err := http.Get(urlStr)
			//defer resp.Body.Close()
			if err != nil {
				fmt.Printf(err.Error())
			}
			// body, err := ioutil.ReadAll(resp.Body)
			// fmt.Printf("response is :%s \n", string(body))
			return resp, err
		}

	}
	// if len(endpointsilices) == 0 {
	// 	resp, err :=http.Get(apiGateway + functionaddress)

	// 	} else {
	// 	for i := range endpointsilices {
	// 		respc, errc := http.Get(endpointsilices[i])
	// 		//defer resp.Body.Close()
	// 		//body, err := ioutil.ReadAll(respc.Body)

	// 	}
	return resp, err
}

func main() {

	functionaddress := "chain-1"

	resp, err := control(functionaddress)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("response is :%s", string(body))
}

// apigatewayaddress := "http://172.16.252.163:31112"

// 	functionName := Get(functionaddress)

// 	fmt.Sprintf("chain-1, input was: %s", functionName)

// 	//var kubeconfig string
// 	var masterURL string

// 	masterURL = "https://172.16.252.163:6443"

// 	log.Printf("Version 0.0 2022-06-29")

// 	var kubeconfig *string
// 	if home := homeDir(); home != "" {
// 		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
// 	} else {
// 		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
// 	}
// 	flag.Parse()

// 	clientCmdConfig, err := clientcmd.BuildConfigFromFlags(masterURL, *kubeconfig)
// 	if err != nil {
// 		log.Fatalf("Error building kubeconfig: %s", err.Error())
// 	}

// 	// flag.StringVar(&kubeconfig, "config", "/home/master1/Downloads/control-d/config",
// 	// 	"Path to a kubeconfig. Only required if out-of-cluster.")

// 	// flag.StringVar(&masterURL, "server", "/home/master1/Downloads/control-d/config",
// 	// 	"The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")

// 	readConfig := config.ReadConfig{}
// 	osEnv := providertypes.OsEnv{}
// 	config, err := readConfig.Read(osEnv)

// 	if err != nil {
// 		log.Fatalf("Error reading config: %s", err.Error())
// 	}

// 	deployConfig := k8s.DeploymentConfig{
// 		RuntimeHTTPPort: 8080,
// 		HTTPProbe:       config.HTTPProbe,
// 		SetNonRootUser:  config.SetNonRootUser,
// 		ReadinessProbe: &k8s.ProbeConfig{
// 			InitialDelaySeconds: int32(config.ReadinessProbeInitialDelaySeconds),
// 			TimeoutSeconds:      int32(config.ReadinessProbeTimeoutSeconds),
// 			PeriodSeconds:       int32(config.ReadinessProbePeriodSeconds),
// 		},
// 		LivenessProbe: &k8s.ProbeConfig{
// 			InitialDelaySeconds: int32(config.LivenessProbeInitialDelaySeconds),
// 			TimeoutSeconds:      int32(config.LivenessProbeTimeoutSeconds),
// 			PeriodSeconds:       int32(config.LivenessProbePeriodSeconds),
// 		},
// 		ImagePullPolicy:   config.ImagePullPolicy,
// 		ProfilesNamespace: config.ProfilesNamespace,
// 	}

// 	kubeconfigQPS := 100
// 	kubeconfigBurst := 250

// 	clientCmdConfig.QPS = float32(kubeconfigQPS)
// 	clientCmdConfig.Burst = kubeconfigBurst

// 	kubeClient, err := kubernetes.NewForConfig(clientCmdConfig)
// 	if err != nil {
// 		log.Fatalf("Error building Kubernetes clientset: %s", err.Error())
// 	}

// 	// pod, err := kubeClient.CoreV1().Pods("keda").Get(context.TODO(), "keda-operator-6475c64bc-mtn8z", metav1.GetOptions{})
// 	// if err != nil {
// 	// 	panic(err.Error())
// 	// }
// 	// fmt.Printf("%v\n\n\n\n", pod.Spec)
// 	//deployment, err := clientset.AppsV1beta1().Deployments("keda").Get("keda-operator", metav1.GetOptions{})
// 	pods, err := kubeClient.CoreV1().Pods(config.DefaultFunctionNamespace).List(context.TODO(), metav1.ListOptions{})
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	for _, v := range pods.Items {
// 		if strings.Contains(v.Name, "test") {
// 			//functionName = strings.TrimSuffix(name, "."+namespace)
// 			fmt.Printf("Name: %v IP: %v\n", v.Name, v.Status.PodIP)
// 		}
// 		//fmt.Printf("Name: %v IP: %v\n", v.Name, v.Status.PodIP)
// 	}

// 	faasClient, err := clientset.NewForConfig(clientCmdConfig)
// 	if err != nil {
// 		log.Fatalf("Error building OpenFaaS clientset: %s", err.Error())
// 	}

// 	defaultResync := time.Minute * 5

// 	namespaceScope := config.DefaultFunctionNamespace
// 	if config.ClusterRole {
// 		namespaceScope = ""
// 	}

// 	kubeInformerOpt := kubeinformers.WithNamespace(namespaceScope)
// 	kubeInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(kubeClient, defaultResync, kubeInformerOpt)

// 	faasInformerOpt := informers.WithNamespace(namespaceScope)
// 	faasInformerFactory := informers.NewSharedInformerFactoryWithOptions(faasClient, defaultResync, faasInformerOpt)

// 	// this is where we need to swap to the faasInformerFactory
// 	profileInformerOpt := informers.WithNamespace(config.ProfilesNamespace)
// 	profileInformerFactory := informers.NewSharedInformerFactoryWithOptions(faasClient, defaultResync, profileInformerOpt)

// 	profileLister := profileInformerFactory.Openfaas().V1().Profiles().Lister()
// 	factory := k8s.NewFunctionFactory(kubeClient, deployConfig, profileLister)

// 	setup := serverSetup{
// 		config:                 config,
// 		functionFactory:        factory,
// 		kubeInformerFactory:    kubeInformerFactory,
// 		faasInformerFactory:    faasInformerFactory,
// 		profileInformerFactory: profileInformerFactory,
// 		kubeClient:             kubeClient,
// 		faasClient:             faasClient,
// 	}

// 	stopCh := signals.SetupSignalHandler()
// 	operator := false
// 	listers := startInformers(setup, stopCh, operator)
// 	functionLookup := NewFunctionLookup(config.DefaultFunctionNamespace, listers.EndpointsInformer.Lister())

// 	functionAddr, resolveErr := functionLookup.Resolve(functionName)
// 	// kaiyu
// 	log.Printf("FunctionName: %s, ResolveAddr: %s", functionName, functionAddr)

// 	if resolveErr != nil {
// 		// TODO: Should record the 404/not found error in Prometheus.
// 		log.Printf("resolver error: no endpoints for %s: %s\n", functionName, resolveErr.Error())
// 		return
// 	}

// 	// resp2, err := http.Get(apigatewayaddress)

// }

// func Get(url string) string {

// 	name := getServiceName(url)

// 	return name
// }

// func getServiceName(urlValue string) string {
// 	var serviceName string
// 	forward := "/function/"
// 	if strings.HasPrefix(urlValue, forward) {
// 		// With a path like `/function/xyz/rest/of/path?q=a`, the service
// 		// name we wish to locate is just the `xyz` portion.  With a positive
// 		// match on the regex below, it will return a three-element slice.
// 		// The item at index `0` is the same as `urlValue`, at `1`
// 		// will be the service name we need, and at `2` the rest of the path.
// 		matcher := functionMatcher.Copy()
// 		matches := matcher.FindStringSubmatch(urlValue)
// 		if len(matches) == hasPathCount {
// 			serviceName = matches[nameIndex]
// 		}
// 	}
// 	return strings.Trim(serviceName, "/")
// }

// func HasPrefix(s, prefix string) bool {
// 	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
// }

// func startInformers(setup serverSetup, stopCh <-chan struct{}, operator bool) customInformers {
// 	kubeInformerFactory := setup.kubeInformerFactory
// 	faasInformerFactory := setup.faasInformerFactory

// 	var functions v1.FunctionInformer
// 	if operator {
// 		// go faasInformerFactory.Start(stopCh)

// 		functions = faasInformerFactory.Openfaas().V1().Functions()
// 		go functions.Informer().Run(stopCh)
// 		if ok := cache.WaitForNamedCacheSync("faas-netes:functions", stopCh, functions.Informer().HasSynced); !ok {
// 			log.Fatalf("failed to wait for cache to sync")
// 		}
// 	}

// 	// go kubeInformerFactory.Start(stopCh)

// 	deployments := kubeInformerFactory.Apps().V1().Deployments()
// 	go deployments.Informer().Run(stopCh)
// 	if ok := cache.WaitForNamedCacheSync("faas-netes:deployments", stopCh, deployments.Informer().HasSynced); !ok {
// 		log.Fatalf("failed to wait for cache to sync")
// 	}

// 	endpoints := kubeInformerFactory.Core().V1().Endpoints()
// 	go endpoints.Informer().Run(stopCh)
// 	if ok := cache.WaitForNamedCacheSync("faas-netes:endpoints", stopCh, endpoints.Informer().HasSynced); !ok {
// 		log.Fatalf("failed to wait for cache to sync")
// 	}

// 	// go setup.profileInformerFactory.Start(stopCh)

// 	profileInformerFactory := setup.profileInformerFactory
// 	profiles := profileInformerFactory.Openfaas().V1().Profiles()
// 	go profiles.Informer().Run(stopCh)
// 	if ok := cache.WaitForNamedCacheSync("faas-netes:profiles", stopCh, profiles.Informer().HasSynced); !ok {
// 		log.Fatalf("failed to wait for cache to sync")
// 	}

// 	return customInformers{
// 		EndpointsInformer:  endpoints,
// 		DeploymentInformer: deployments,
// 		FunctionsInformer:  functions,
// 	}
// }

// type serverSetup struct {
// 	config                 config.BootstrapConfig
// 	kubeClient             *kubernetes.Clientset
// 	faasClient             *clientset.Clientset
// 	functionFactory        k8s.FunctionFactory
// 	kubeInformerFactory    kubeinformers.SharedInformerFactory
// 	faasInformerFactory    informers.SharedInformerFactory
// 	profileInformerFactory informers.SharedInformerFactory
// }

// type customInformers struct {
// 	EndpointsInformer  v1core.EndpointsInformer
// 	DeploymentInformer v1apps.DeploymentInformer
// 	FunctionsInformer  v1.FunctionInformer
// }

// func NewFunctionLookup(ns string, lister corelister.EndpointsLister) *FunctionLookup {
// 	return &FunctionLookup{
// 		DefaultNamespace: ns,
// 		EndpointLister:   lister,
// 		Listers:          map[string]corelister.EndpointsNamespaceLister{},
// 		lock:             sync.RWMutex{},
// 	}
// }

// type FunctionLookup struct {
// 	DefaultNamespace string
// 	EndpointLister   corelister.EndpointsLister
// 	Listers          map[string]corelister.EndpointsNamespaceLister

// 	lock sync.RWMutex
// }

// func (l *FunctionLookup) Resolve(name string) (url.URL, error) {
// 	functionName := name
// 	namespace := getNamespace(name, l.DefaultNamespace)
// 	if err := l.verifyNamespace(namespace); err != nil {
// 		return url.URL{}, err
// 	}

// 	if strings.Contains(name, ".") {
// 		functionName = strings.TrimSuffix(name, "."+namespace)
// 	}

// 	nsEndpointLister := l.GetLister(namespace)

// 	if nsEndpointLister == nil {
// 		l.SetLister(namespace, l.EndpointLister.Endpoints(namespace))

// 		nsEndpointLister = l.GetLister(namespace)
// 	}

// 	svc, err := nsEndpointLister.Get(functionName)
// 	if err != nil {
// 		return url.URL{}, fmt.Errorf("error listing \"%s.%s\": %s", functionName, namespace, err.Error())
// 	}

// 	if len(svc.Subsets) == 0 {
// 		return url.URL{}, fmt.Errorf("no subsets available for \"%s.%s\"", functionName, namespace)
// 	}

// 	all := len(svc.Subsets[0].Addresses)
// 	if len(svc.Subsets[0].Addresses) == 0 {
// 		return url.URL{}, fmt.Errorf("no addresses in subset for \"%s.%s\"", functionName, namespace)
// 	}

// 	target := rand.Intn(all)

// 	serviceIP := svc.Subsets[0].Addresses[target].IP

// 	urlStr := fmt.Sprintf("http://%s:%d", serviceIP, watchdogPort)

// 	urlRes, err := url.Parse(urlStr)
// 	if err != nil {
// 		return url.URL{}, err
// 	}

// 	log.Printf("[Call k8s/proxy.go Resolve] name: %s, url %s", name, urlStr)

// 	return *urlRes, nil
// }

// func getNamespace(name, defaultNamespace string) string {
// 	namespace := defaultNamespace
// 	if strings.Contains(name, ".") {
// 		namespace = name[strings.LastIndexAny(name, ".")+1:]
// 	}
// 	return namespace
// }

// func (l *FunctionLookup) verifyNamespace(name string) error {
// 	if name != "kube-system" {
// 		return nil
// 	}
// 	// ToDo use global namepace parse and validation
// 	return fmt.Errorf("namespace not allowed")
// }

// func (f *FunctionLookup) GetLister(ns string) corelister.EndpointsNamespaceLister {
// 	f.lock.RLock()
// 	defer f.lock.RUnlock()
// 	return f.Listers[ns]
// }

// func (f *FunctionLookup) SetLister(ns string, lister corelister.EndpointsNamespaceLister) {
// 	f.lock.Lock()
// 	defer f.lock.Unlock()
// 	f.Listers[ns] = lister
// }
