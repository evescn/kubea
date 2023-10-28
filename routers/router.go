package routers

import (
	"github.com/gin-gonic/gin"
	"kubea/controller"
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
		//登录验证，路由权限信息
		POST("/api/login", controller.Login.Auth).
		// 用户管理
		GET("/api/user/list", controller.User.List).
		POST("/api/user/add", controller.User.Add).
		PUT("/api/user/update", controller.User.Update).
		PUT("/api/user/updateAdmin", controller.User.UpdateAdmin).
		PUT("/api/user/updateRole", controller.User.UpdateRole).
		DELETE("/api/user/del", controller.User.Delete).
		// 环境管理
		GET("/api/url/env/list", controller.Env.List).
		POST("/api/url/env/add", controller.Env.Add).
		PUT("/api/url/env/update", controller.Env.Update).
		DELETE("/api/url/env/del", controller.Env.Delete).
		// URL信息管理
		GET("/api/url/svc/list", controller.Service.List).
		POST("/api/url/svc/add", controller.Service.Add).
		PUT("/api/url/svc/update", controller.Service.Update).
		DELETE("/api/url/svc/del", controller.Service.Delete).
		// 角色管理
		GET("/api/role/list", controller.Roles.List).
		GET("/api/role/getAll", controller.Roles.GetAll).
		POST("/api/role/add", controller.Roles.Add).
		PUT("/api/role/update", controller.Roles.Update).
		DELETE("/api/role/del", controller.Roles.Delete).
		// 1级菜单管理
		GET("/api/menu/list", controller.Menu.List).
		POST("/api/menu/add", controller.Menu.Add).
		PUT("/api/menu/update", controller.Menu.Update).
		DELETE("/api/menu/del", controller.Menu.Delete).
		// 2级菜单管理
		GET("/api/submenu/list", controller.SubMenu.List).
		POST("/api/submenu/add", controller.SubMenu.Add).
		PUT("/api/submenu/update", controller.SubMenu.Update).
		DELETE("/api/submenu/del", controller.SubMenu.Delete).
		// 3级菜单管理
		GET("/api/subsubmenu/list", controller.SubSubMenus.List).
		POST("/api/subsubmenu/add", controller.SubSubMenus.Add).
		PUT("/api/subsubmenu/update", controller.SubSubMenus.Update).
		DELETE("/api/subsubmenu/del", controller.SubSubMenus.Delete).
		// 权限菜单关系管理
		GET("/api/roleMenuRelation/getAll", controller.RoleMenuRelation.GetAll).
		GET("/api/roleMenuRelation/getPermissions", controller.RoleMenuRelation.GetPermissions).
		PUT("/api/roleMenuRelation/update", controller.RoleMenuRelation.Update).
		//应用管理
		GET("/api/app/list", controller.App.List).
		GET("/api/app/get", controller.App.Get).
		GET("/api/app/all", controller.App.GetAll).
		POST("/api/app/add", controller.App.Add).
		PUT("/api/app/update", controller.App.Update).
		DELETE("/api/app/del", controller.App.Delete).
		GET("/api/app/getApp", controller.App.GetApp).
		//部署记录
		GET("/api/deploy/list", controller.Deploy.List).
		POST("/api/deploy/add", controller.Deploy.Add).
		PUT("/api/deploy/update", controller.Deploy.Update).
		DELETE("/api/deploy/del", controller.Deploy.Delete).
		GET("/api/deploy/getLog", controller.Deploy.GetLog).
		POST("/api/deploy/cicd", controller.Deploy.CiCd).
		POST("/api/deploy/jenkniscicd", controller.Deploy.JenkinsCiCd).
		POST("/api/deploy/updatecicd", controller.Deploy.UpdateCiCd).
		//集群
		GET("/api/k8s/clusters", controller.Cluster.GetClusters).
		// Pod 操作
		GET("/api/k8s/pods", controller.Pod.GetPods).
		GET("/api/k8s/pod/detail", controller.Pod.GetPodDetail).
		DELETE("/api/k8s/pod", controller.Pod.DeletePod).
		PUT("/api/k8s/pod", controller.Pod.UpdatePod).
		GET("/api/k8s/pod/container", controller.Pod.GetPodContainer).
		GET("/api/k8s/pod/log", controller.Pod.GetPodLog).
		// Deployment 操作
		GET("/api/k8s/deployments", controller.Deployment.GetDeployments).
		GET("/api/k8s/deployment/detail", controller.Deployment.GetDeploymentDetail).
		DELETE("/api/k8s/deployment", controller.Deployment.DeleteDeployment).
		PUT("/api/k8s/deployment", controller.Deployment.UpdateDeployment).
		PUT("/api/k8s/deployment/scale", controller.Deployment.ScaleDeployment).
		PUT("/api/k8s/deployment/restart", controller.Deployment.RestartDeployment).
		POST("/api/k8s/deployment/create", controller.Deployment.CreateDeployment).
		// DaemonSet 操作
		GET("/api/k8s/daemonsets", controller.DaemonSet.GetDaemonSets).
		GET("/api/k8s/daemonset/detail", controller.DaemonSet.GetDaemonSetDetail).
		DELETE("/api/k8s/daemonset", controller.DaemonSet.DeleteDaemonSet).
		PUT("/api/k8s/daemonset", controller.DaemonSet.UpdateDaemonSet).
		// StatefulSet 操作
		GET("/api/k8s/statefulsets", controller.StatefulSet.GetStatefulSets).
		GET("/api/k8s/statefulset/detail", controller.StatefulSet.GetStatefulSetDetail).
		DELETE("/api/k8s/statefulset", controller.StatefulSet.DeleteStatefulSet).
		PUT("/api/k8s/statefulset", controller.StatefulSet.UpdateStatefulSet).
		// Service 操作
		GET("/api/k8s/services", controller.Servicev1.GetServices).
		GET("/api/k8s/service/detail", controller.Servicev1.GetServiceDetail).
		DELETE("/api/k8s/service", controller.Servicev1.DeleteService).
		PUT("/api/k8s/service", controller.Servicev1.UpdateService).
		POST("/api/k8s/service/create", controller.Servicev1.CreateService).
		// Ingress 操作
		GET("/api/k8s/ingresses", controller.Ingress.GetIngresses).
		GET("/api/k8s/ingress/detail", controller.Ingress.GetIngressDetail).
		DELETE("/api/k8s/ingress", controller.Ingress.DeleteIngress).
		PUT("/api/k8s/ingress", controller.Ingress.UpdateIngress).
		POST("/api/k8s/ingress/create", controller.Ingress.CreateIngress).
		// Node 操作
		GET("/api/k8s/nodes", controller.Node.GetNodes).
		GET("/api/k8s/node/detail", controller.Node.GetNodeDetail).
		// Namespace 操作
		GET("/api/k8s/namespaces", controller.Namespace.GetNamespaces).
		GET("/api/k8s/namespace/detail", controller.Namespace.GetNamespaceDetail).
		DELETE("/api/k8s/namespace", controller.Namespace.DeleteNamespace).
		// PV 操作
		GET("/api/k8s/pvs", controller.Pv.GetPvs).
		GET("/api/k8s/pv/detail", controller.Pv.GetPvDetail).
		DELETE("/api/k8s/pv", controller.Pv.DeletePv).
		// ConfigMap 操作
		GET("/api/k8s/configmaps", controller.ConfigMap.GetConfigMaps).
		GET("/api/k8s/configmap/detail", controller.ConfigMap.GetConfigMapDetail).
		DELETE("/api/k8s/configmap", controller.ConfigMap.DeleteConfigMap).
		PUT("/api/k8s/configmap", controller.ConfigMap.UpdateConfigMap).
		// Secret 操作
		GET("/api/k8s/secrets", controller.Secret.GetSecrets).
		GET("/api/k8s/secret/detail", controller.Secret.GetSecretDetail).
		DELETE("/api/k8s/secret", controller.Secret.DeleteSecret).
		PUT("/api/k8s/secret", controller.Secret.UpdateSecret).
		// PVC 操作
		GET("/api/k8s/pvcs", controller.Pvc.GetPvcs).
		GET("/api/k8s/pvc/detail", controller.Pvc.GetPvcDetail).
		DELETE("/api/k8s/pvc", controller.Pvc.DeletePvc).
		PUT("/api/k8s/pvc", controller.Pvc.UpdatePvc).
		// Event 操作
		GET("/api/k8s/events", controller.Event.GetList).
		// AllRes 操作
		GET("/api/k8s/allres", controller.AllRes.GetAllNum).
		// Helm 应用商店
		GET("/api/helmstore/releases", controller.HelmStore.ListReleases).
		GET("/api/helmstore/release/detail", controller.HelmStore.DetailRelease).
		POST("/api/helmstore/release/install", controller.HelmStore.InstallRelease).
		DELETE("/api/helmstore/release/uninstall", controller.HelmStore.UninstallRelease).
		GET("/api/helmstore/charts", controller.HelmStore.ListCharts).
		POST("/api/helmstore/chart/add", controller.HelmStore.AddChart).
		PUT("/api/helmstore/chart/update", controller.HelmStore.UpdateChart).
		DELETE("/api/helmstore/chart/del", controller.HelmStore.DeleteChart).
		POST("/api/helmstore/chartfile/upload", controller.HelmStore.UploadChartFile).
		DELETE("/api/helmstore/chartfile/del", controller.HelmStore.DeleteChartFile)
}
