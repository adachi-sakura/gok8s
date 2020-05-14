package algorithm

import (
	"context"
	"github.com/buzaiguna/gok8s/appctx"
	"github.com/buzaiguna/gok8s/utils"
	"testing"
)

type (
	testData struct {
		ctx 	context.Context
		query 	map[string]string
	}
)


func dummyData() testData {
	yamlFiles := `  apiVersion: apps/v1
  kind: Deployment
  metadata:
    creationTimestamp: null
    labels:
      name: carts-db
    annotations:
      leastResponseTime: '0.5'
      httpRequestCount: '48000'
    name: carts-db
    namespace: sock-shop
  spec:
    progressDeadlineSeconds: 2147483647
    replicas: 1
    revisionHistoryLimit: 2147483647
    selector:
      matchLabels:
        name: carts-db
    strategy:
      rollingUpdate:
        maxSurge: 1
        maxUnavailable: 1
      type: RollingUpdate
    template:
      metadata:
        creationTimestamp: null
        labels:
          name: carts-db
      spec:
        containers:
        - image: mongo
          imagePullPolicy: Always
          name: carts-db
          ports:
          - containerPort: 27017
            name: mongo
            protocol: TCP
          resources: {}
          securityContext:
            capabilities:
              add:
              - CHOWN
              - SETGID
              - SETUID
              drop:
              - all
            readOnlyRootFilesystem: true
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
          - mountPath: /tmp
            name: tmp-volume
        dnsPolicy: ClusterFirst
        nodeSelector:
          beta.kubernetes.io/os: linux
        restartPolicy: Always
        schedulerName: default-scheduler
        securityContext: {}
        terminationGracePeriodSeconds: 30
        volumes:
        - emptyDir:
            medium: Memory
          name: tmp-volume
  status: {}
---
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    creationTimestamp: null
    labels:
      name: carts
    name: carts
    annotations:
      leastResponseTime: '0.8'
      dependencies: 'carts-db'
      httpRequestCount: '48000'
    namespace: sock-shop
  spec:
    progressDeadlineSeconds: 2147483647
    replicas: 1
    revisionHistoryLimit: 2147483647
    selector:
      matchLabels:
        name: carts
    strategy:
      rollingUpdate:
        maxSurge: 1
        maxUnavailable: 1
      type: RollingUpdate
    template:
      metadata:
        creationTimestamp: null
        labels:
          name: carts
      spec:
        containers:
        - env:
          - name: ZIPKIN
            value: zipkin.jaeger.svc.cluster.local
          - name: JAVA_OPTS
            value: -Xms64m -Xmx128m -XX:PermSize=32m -XX:MaxPermSize=64m -XX:+UseG1GC
              -Djava.security.egd=file:/dev/urandom
          image: weaveworksdemos/carts:0.4.8
          imagePullPolicy: IfNotPresent
          name: carts
          ports:
          - containerPort: 80
            protocol: TCP
          resources: {}
          securityContext:
            capabilities:
              add:
              - NET_BIND_SERVICE
              drop:
              - all
            readOnlyRootFilesystem: true
            runAsNonRoot: true
            runAsUser: 10001
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
          - mountPath: /tmp
            name: tmp-volume
        dnsPolicy: ClusterFirst
        nodeSelector:
          beta.kubernetes.io/os: linux
        restartPolicy: Always
        schedulerName: default-scheduler
        securityContext: {}
        terminationGracePeriodSeconds: 30
        volumes:
        - emptyDir:
            medium: Memory
          name: tmp-volume
  status: {}`
	objects := utils.ParseK8SYaml([]byte(yamlFiles))
	ctx := context.Background()
	ctx = appctx.WithK8SObjects(ctx, objects)
	deployments := appctx.DeploymentObjects(ctx)
	ctx = appctx.WithDeployments(ctx, deployments)
	ctx = appctx.DeploymentInvertedIndexContext(ctx, deployments)

	query := map[string]string {
		"entry":	"carts",
		"totalTime":	"2",
	}


	return testData{ctx, query}
}


func TestValidate(t *testing.T) {
	ctx := dummyData().ctx
	if err := validate(ctx); err != nil {
		t.Error("validate test failed")
	}
}

func TestBuildMicroserviceYaml(t *testing.T) {
	ctx := dummyData().ctx
	deployments := appctx.Deployments(ctx)
	for _, deployment := range deployments {
		if _, err := buildMicroserviceYaml(ctx, deployment); err != nil {
			t.Error("microservice yaml test failed")
		}
	}

}

func TestBuildEntrance(t *testing.T) {
	data := dummyData()
	ctx := data.ctx
	q := data.query
	if _, err := appctx.GetDeploymentIndex(ctx, q["entry"]); err != nil {
		t.Error("build entrance test failed")
	}

}