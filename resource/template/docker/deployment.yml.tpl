apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    k8s-app: eztalk-monitor
    qcloud-app: eztalk-monitor
  name: eztalk-monitor
  namespace: application-services
spec:
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      k8s-app: eztalk-monitor
      qcloud-app: eztalk-monitor
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 0
    type: RollingUpdate
  template:
    metadata:
      labels:
        k8s-app: eztalk-monitor
        qcloud-app: eztalk-monitor
    spec:
      imagePullSecrets:
        - name: mowen-images
      restartPolicy: Always
      containers:
        - name: eztalk-monitor
          image: qdtech-vpc.tencentcloudcr.com/application-services/eztalk-monitor:${trigger["artifacts"][0]["version"]}
          imagePullPolicy: IfNotPresent
          readinessProbe:
            httpGet:
              scheme: HTTP
              path: /ping
              port: 8080
            initialDelaySeconds: 20
            periodSeconds: 2
            successThreshold: 5
            failureThreshold: 3
          ports:
            - name: http
              containerPort: 8080
          resources:
            limits:
              cpu: 1000m
              memory: 1Gi
            requests:
              cpu: 500m
              memory: 256Mi
          env:
            - name: APP_NAME
              value: eztalk-monitor
            - name: APP_ENV
              valueFrom:
                configMapKeyRef:
                  name: env
                  key: APP_ENV
            - name: TIMEZONE
              valueFrom:
                configMapKeyRef:
                  name: env
                  key: TIMEZONE
            - name: NACOS_HOST
              valueFrom:
                configMapKeyRef:
                  name: nacos-cConfig
                  key: host
            - name: NACOS_PORT
              valueFrom:
                configMapKeyRef:
                  name: nacos-cConfig
                  key: port
            - name: NACOS_USERNAME
              valueFrom:
                configMapKeyRef:
                  name: nacos-cConfig
                  key: username
            - name: NACOS_PASSWORD
              valueFrom:
                configMapKeyRef:
                  name: nacos-cConfig
                  key: password