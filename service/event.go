package service

import (
	"fmt"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"kubea/dao"
	"kubea/model"
	"time"
)

var Event event

type event struct{}

// GetList 获取列表
func (*event) GetList(name, cluster string, page, limit int) (*dao.Events, error) {
	data, err := dao.Event.GetList(name, cluster, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WatchEventTask informer监听
func (*event) WatchEventTask(cluster string) {
	// 实例化informer
	informerFactory := informers.NewSharedInformerFactory(K8s.ClientMap[cluster], time.Minute)
	// 监听资源
	informer := informerFactory.Core().V1().Events()
	// 添加事件handler
	informer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				onAdd(obj, cluster)
			},
		},
	)

	// 处理启动和优雅关闭
	stopCh := make(chan struct{})
	defer close(stopCh)
	informerFactory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, informer.Informer().HasSynced) {
		logger.Error("同步cache超时")
		return
	}
	<-stopCh

}

// onAdd 新增落库
func onAdd(obj interface{}, cluster string) {
	// 断言
	event := obj.(*corev1.Event)
	// 判断是否重复 name, kind, namespace, reason string, eventTime time.Time, cluster string
	_, has, err := dao.Event.HasEvent(event.InvolvedObject.Name,
		event.InvolvedObject.Kind,
		event.InvolvedObject.Namespace,
		event.Reason,
		event.CreationTimestamp.Time,
		cluster,
	)
	if err != nil {
		return
	}
	if has {
		logger.Info(fmt.Sprintf("Event数据已存在, %s %s %s %s %v %s",
			event.InvolvedObject.Name,
			event.InvolvedObject.Kind,
			event.InvolvedObject.Namespace,
			event.Reason,
			event.CreationTimestamp.Time,
			cluster),
		)
	}
	// 组装数据
	data := &model.Event{
		Name:      event.InvolvedObject.Name,
		Kind:      event.InvolvedObject.Kind,
		Namespace: event.InvolvedObject.Namespace,
		Rtype:     event.Type,
		Reason:    event.Reason,
		Message:   event.Message,
		EventTime: &event.CreationTimestamp.Time,
		Cluster:   cluster,
	}
	// 数据库添加
	if err := dao.Event.Add(data); err != nil {
		return
	}

}
