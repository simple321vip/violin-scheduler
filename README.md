## violin-scheduler

[参照文章，十分感谢](https://www.qikqiak.com/post/custom-kube-scheduler/)

### 配置文件挂载

--config 指定的是 /etc/kubernetes/scheduler-extender.yaml 
鉴于云服务可能不会将权限提供给客户，所以采用configMap形式挂载 自定义调度器配置文件

    // /etc/kubernetes/scheduler-extender.yaml
    --- 
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: violin-scheduler-config
      namespace: devops
    data:
      application: |-
        apiVersion: kubescheduler.config.k8s.io/v1alpha1
        kind: KubeSchedulerConfiguration
        clientConnection:
          kubeconfig: "/etc/kubernetes/scheduler.conf"
        algorithmSource:
          policy:
            file:
              path: "/etc/kubernetes/scheduler-extender-policy.yaml"
    
而自定义调度策略文件 /etc/kubernetes/scheduler-extender-policy.yaml 我们也是基于yaml来写的

    // /etc/kubernetes/scheduler-extender-policy.yaml
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: violin-scheduler-config
      namespace: devops
    data:
      app: |-
        apiVersion: v1
        kind: Policy
        extenders:
        - urlPrefix: "http://127.0.0.1:8888/"
          filterVerb: "filter"
          prioritizeVerb: "prioritize"
          weight: 1
          enableHttps: false

## 
 
 发布多个调度器 下面文件
 /etc/kubernetes/manifests/kube-schduler.yaml