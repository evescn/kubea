package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var Router router

type router struct{}

func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/testapi", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "testapi success!",
			"data": nil,
		})
	}).
		//登录验证
		POST("/api/login", Login.Auth).
		//应用管理
		GET("/api/app/list", App.List).
		GET("/api/app/all", App.GetAll).
		POST("/api/app/add", App.Add).
		PUT("/api/app/update", App.Update).
		DELETE("/api/app/del", App.Delete).
		//GET("/api/app/alltags", App.GetAllTags).
		//编排管理
		GET("/api/orchestration/list", Orchestration.List).
		POST("/api/orchestration/add", Orchestration.Add).
		PUT("/api/orchestration/update", Orchestration.Update).
		DELETE("/api/orchestration/del", Orchestration.Delete).
		//集群
		GET("/api/k8s/clusters", Cluster.GetClusters).
		// Pod 操作
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod", Pod.DeletePod).
		PUT("/api/k8s/pod", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodLog).
		// Deployment 操作
		GET("/api/k8s/deployments", Deployment.GetDeployments).
		GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail).
		DELETE("/api/k8s/deployment", Deployment.DeleteDeployment).
		PUT("/api/k8s/deployment", Deployment.UpdateDeployment).
		PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment).
		PUT("/api/k8s/deployment/restart", Deployment.RestartDeployment).
		POST("/api/k8s/deployment/create", Deployment.CreateDeployment).
		// DaemonSet 操作
		GET("/api/k8s/daemonsets", DaemonSet.GetDaemonSets).
		GET("/api/k8s/daemonset/detail", DaemonSet.GetDaemonSetDetail).
		DELETE("/api/k8s/daemonset", DaemonSet.DeleteDaemonSet).
		PUT("/api/k8s/daemonset", DaemonSet.UpdateDaemonSet).
		// StatefulSet 操作
		GET("/api/k8s/statefulsets", StatefulSet.GetStatefulSets).
		GET("/api/k8s/statefulset/detail", StatefulSet.GetStatefulSetDetail).
		DELETE("/api/k8s/statefulset", StatefulSet.DeleteStatefulSet).
		PUT("/api/k8s/statefulset", StatefulSet.UpdateStatefulSet).
		// Service 操作
		GET("/api/k8s/services", Servicev1.GetServices).
		GET("/api/k8s/service/detail", Servicev1.GetServiceDetail).
		DELETE("/api/k8s/service", Servicev1.DeleteService).
		PUT("/api/k8s/service", Servicev1.UpdateService).
		POST("/api/k8s/service/create", Servicev1.CreateService).
		// Ingress 操作
		GET("/api/k8s/ingresses", Ingress.GetIngresses).
		GET("/api/k8s/ingress/detail", Ingress.GetIngressDetail).
		DELETE("/api/k8s/ingress", Ingress.DeleteIngress).
		PUT("/api/k8s/ingress", Ingress.UpdateIngress).
		POST("/api/k8s/ingress/create", Ingress.CreateIngress).
		// Node 操作
		GET("/api/k8s/nodes", Node.GetNodes).
		GET("/api/k8s/node/detail", Node.GetNodeDetail).
		// Namespace 操作
		GET("/api/k8s/namespaces", Namespace.GetNamespaces).
		GET("/api/k8s/namespace/detail", Namespace.GetNamespaceDetail).
		DELETE("/api/k8s/namespace", Namespace.DeleteNamespace).
		// PV 操作
		GET("/api/k8s/pvs", Pv.GetPvs).
		GET("/api/k8s/pv/detail", Pv.GetPvDetail).
		DELETE("/api/k8s/pv", Pv.DeletePv).
		// ConfigMap 操作
		GET("/api/k8s/configmaps", ConfigMap.GetConfigMaps).
		GET("/api/k8s/configmap/detail", ConfigMap.GetConfigMapDetail).
		DELETE("/api/k8s/configmap", ConfigMap.DeleteConfigMap).
		PUT("/api/k8s/configmap", ConfigMap.UpdateConfigMap).
		// Secret 操作
		GET("/api/k8s/secrets", Secret.GetSecrets).
		GET("/api/k8s/secret/detail", Secret.GetSecretDetail).
		DELETE("/api/k8s/secret", Secret.DeleteSecret).
		PUT("/api/k8s/secret", Secret.UpdateSecret).
		// PVC 操作
		GET("/api/k8s/pvcs", Pvc.GetPvcs).
		GET("/api/k8s/pvc/detail", Pvc.GetPvcDetail).
		DELETE("/api/k8s/pvc", Pvc.DeletePvc).
		PUT("/api/k8s/pvc", Pvc.UpdatePvc).
		// Event 操作
		GET("/api/k8s/events", Event.GetList).
		// AllRes 操作
		GET("/api/k8s/allres", AllRes.GetAllNum).
		// Helm 应用商店
		GET("/api/helmstore/releases", HelmStore.ListReleases).
		GET("/api/helmstore/release/detail", HelmStore.DetailRelease).
		POST("/api/helmstore/release/install", HelmStore.InstallRelease).
		DELETE("/api/helmstore/release/uninstall", HelmStore.UninstallRelease).
		GET("/api/helmstore/charts", HelmStore.ListCharts).
		POST("/api/helmstore/chart/add", HelmStore.AddChart).
		PUT("/api/helmstore/chart/update", HelmStore.UpdateChart).
		DELETE("/api/helmstore/chart/del", HelmStore.DeleteChart).
		POST("/api/helmstore/chartfile/upload", HelmStore.UploadChartFile).
		DELETE("/api/helmstore/chartfile/del", HelmStore.DeleteChartFile)
}
