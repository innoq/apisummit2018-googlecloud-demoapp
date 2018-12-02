apiVersion: apps/v1
kind: Deployment
metadata:
  name: todo-deployment
  labels:
    app: todo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: todo
  template:
    metadata:
      labels:
        app: todo
    spec:
      containers:
      - name: todo
        image: gcr.io/PROJECT_ID/todo
        command:
          - ./todo 
          - -addr 
          - 0.0.0.0:8080
          #- -db
          #- postgres://todo:todo1234@35.246.186.44/todo?sslmode=disable
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: todo
spec:
  type: NodePort
  selector:
    app: todo
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: todo
spec:
  rules:
  - http:
      paths:
      - path: /*
        backend:
          serviceName: todo
          servicePort: 80
