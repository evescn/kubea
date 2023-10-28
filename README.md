# kubea

> k8s+cicd+内网URL导航后端

## 1. 服务镜像打包

### 克隆代码

```shell
$ git clone https://github.com/evescn/kubea.git
$ cd kubea
$ git checkout main
```

### 打包镜像

> 方法1 打包 Docker 镜像

```shell
# 第一种打包 Docker 镜像
# 编译项目
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# 打包 Docker 镜像
$ docker build -t harbor.xxx.cn/devops/kubea:v1.1 -f Dockerfile .
$ docker push harbor.xxx.cn/devops/kubea:v1.1
```

> 方法2 打包 Docker 镜像

```shell
# 第二种打包 Docker 镜像
$ chmod a+x ./build
$ ./build 1 # 版本号信息
```

## 2. 服务部署

### a | Docker 启动

```shell
$ docker run -d \
  --restart=always \
  --name kubea \
  -p 8080:8080 \
  harbor.xxx.cn/devops/kubea:v1.1
```

### b | Kubernetes 启动

```shell
# k8s.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: site-navigation
  namespace: devops
spec:
  replicas: 1
  selector:
    matchLabels:
      app: site-navigation
  template:
    metadata:
      labels:
        app: site-navigation
    spec:
      containers:
      - name: site-navigation
        image: harbor.xxx.cn/devops/kubea:v1.0
        imagePullPolicy: Always
        ports:
        - containerPort: 8080

---
# service
apiVersion: v1
kind: Service
metadata:
  name: site-navigation
  namespace: devops
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
    nodePort: 30080
  selector:
    app: site-navigation
  type: NodePort
```

```shell
$ kubectl apply -f k8s.yaml
```

## 3. 服务访问

> 项目前后端分离，需要部署前端后才能访问
> [前端地址](https://github.com/evescn/kubea-fe)


