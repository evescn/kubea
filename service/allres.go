package service

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
)

var AllRes allRes

type allRes struct{}

var mt sync.Mutex

func (*allRes) GetAllNum(client *kubernetes.Clientset) (map[string]int, []error) {
	//等待所有的goroutine执行完之后，再往下执行,这里其实是阻塞的作用
	var wg sync.WaitGroup
	wg.Add(12)

	errs := make([]error, 0)
	//map[资源名]资源数量
	data := make(map[string]int, 0)

	go func() {
		list, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}

		//data["Nodes"] = len(list.Items)
		//为什么要封装addMap方法？
		//因为有12个协程会对这个map进行操作，map默认是线程非安全的，也就是所有协程一起操作map时，会有并发的报错
		//同一时间只能有一个协程对map进行读写操作，所以addMap实际上给map加把锁
		addMap(data, "Nodes", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Namespaces", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Ingresses", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "PVs", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.AppsV1().DaemonSets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "DaemonSets", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "StatefulSets", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.BatchV1().Jobs("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Jobs", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Services", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Deployments", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Pods", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "Secrets", len(list.Items))
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		}
		addMap(data, "ConfigMaps", len(list.Items))
		wg.Done()
	}()
	//当wg里面的计数器为0时，就取消阻塞，继续执行，非0时，则一直阻塞
	wg.Wait()
	return data, nil
}

func addMap(mp map[string]int, resource string, num int) {
	mt.Lock()
	defer mt.Unlock()
	mp[resource] = num
}
