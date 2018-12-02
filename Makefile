PROJECT_ID := ${shell gcloud config get-value project}

local: 
	go build

docker-image:
	docker build -t gcr.io/${PROJECT_ID}/todo .
	docker push gcr.io/${PROJECT_ID}/todo
	m4 -DPROJECT_ID=${PROJECT_ID} kubernetes-todo.yaml.m4 > kubernetes-todo.yaml
