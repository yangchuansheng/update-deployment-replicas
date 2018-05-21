# 在 Kubernetes 的 Pod 中对 deployment 进行扩容和缩容

该项目主要是为了实现定时扩容与缩容 `Deployment` 的功能，额外还支持定时修改 `Deployment` 的镜像。

## 参数说明

```bash
$ ./update-deployment-replicas -h

Usage of ./update-deployment-replicas:
  -alsologtostderr
    	log to standard error as well as files
  -deployment string
    	deployment name
  -image string
    	new image name
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -logtostderr
    	log to standard error instead of files
  -namespace string
    	namespace name (default "default")
  -replicas int
    	number of replicas
  -stderrthreshold value
    	logs at or above this threshold go to stderr
  -v value
    	log level for V logs
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging
```

+ 必须指定参数 `-deployment` 的值
+ 如果不指定参数 `-namespace` 的值，默认使用 `default namespace`
+ 参数 `-replicas` 表示 `Deployment` 的副本数。如果你想对 `Deployment` 进行扩容或缩容，可以通过此参数来设定
+ 参数 `-image` 表示 `Deployment` 使用的镜像名。如果你想修改 `Deployment` 的镜像，可以通过此参数来设定

示例：

+ 将命名空间 `default` 下的 Deployment `nginx` 的实例数变成2：

```bash
$ ./update-deployment-replicas -deployment nginx -replicas 2
```

+ 将命名空间 `demo` 下的 Deployment `nginx` 的镜像修改为 `nginx:alpine`：

```bash
$ ./update-deployment-replicas -namespace demo -deployment nginx -image nginx:alpine
```

## 使用方法

假设 Namespace `default` 下的 Deployment `nginx` 原来的实例数为 1。

#### **场景一：** 每天上午 10:00 将 Deployment `nginx` 的实例数扩容到 2 个

创建一个 CronJob：

```yaml
$ cat scale-job.yaml

apiVersion: batch/v2alpha1
kind: CronJob
metadata:
  name: scale-nginx
spec:
  schedule: "0 10 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: scale-nginx
            image: yangchuansheng/update-deployment-replicas:1.7
          - args:
            - --deployment
            - nginx
            - --replicas
            - '2'
          restartPolicy: OnFailure
```

```bash
$ kubectl create -f scale-job.yaml
```

```bash
$ kubectl get cronjob
NAME            SCHEDULE      SUSPEND   ACTIVE    LAST-SCHEDULE
scale-nginx     0 10 * * *    False     0         <none>
```

#### **场景二：** 每天下午 2:00 将 Deployment `nginx` 的示例数缩容到 1 个：

创建一个 CronJob：

```yaml
$ cat down-job.yaml

apiVersion: batch/v2alpha1
kind: CronJob
metadata:
  name: down-nginx
spec:
  schedule: "0 14 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: down-nginx
            image: yangchuansheng/update-deployment-replicas:1.7
          - args:
            - --deployment
            - nginx
            - --replicas
            - '1'
          restartPolicy: OnFailure
```

```bash
$ kubectl create -f down-job.yaml
```

```bash
$ kubectl get cronjob
NAME            SCHEDULE      SUSPEND   ACTIVE    LAST-SCHEDULE
down-nginx      0 14 * * *    False     0         <none>
```


## 版权

Copyright 2018 Ryan (yangchuansheng33@gmail.com)

MIT License，详情见 LICENSE 文件。
