# todo

### Dieses Repository beinhaltet eine kleine golang todo app, Konfigurationsdateien und Beispiele für den Google Cloud Workshop

### Cloud Shell  
  In dieser Demo erkunden wir die Google Cloud Shell, starten eine VM und deployen die Demo app
 * clonen Sie das Repo der todo App: `git clone https://github.com/innoq/apisummit2018-googlecloud-demoapp`
 * starten Sie die Cloud Shell und bauen Sie die Demo mit `go build`. 
 * staten Sie die App mit `./todo`
 * Testen Sie den HTTP-Zugriff mit [Web Preview](https://cloud.google.com/shell/docs/using-web-preview)

### Compute Engine  
  wir testen `./todo` auf einer VM
  * starten Sie eine VM unter [Compute Engine / VM instances](https://console.cloud.google.com/compute/instances)
  * übertragen Sie `./todo` mit `gcloud compute scp`
  * starten Sie auf der VM die App mit `./todo -addr 0.0.0.0:80`
  * Konfigurieren sie die [Firewall](https://console.cloud.google.com/networking/firewalls) für Port 80
  * Testen Sie den HTTP-Zugriff
  
### Kubernetes Engine  
  Wir erstellen ein Docker Image, und installieren dieses in Kubernetes
 * Starten Sie einen [Cluster](https://console.cloud.google.com/kubernetes)
 * machen Sie sich mit den Kubernetes-Konzepten:
  * [Deployment](https://cloud.google.com/kubernetes-engine/docs/concepts/deployment)
  * [Service](https://cloud.google.com/kubernetes-engine/docs/concepts/service)
  * und [Ingress](https://cloud.google.com/kubernetes-engine/docs/concepts/ingress) vertraut
 * mit `gcloud auth configure-docker` verschaffen Sie sich bzw. Docker Zugriff auf die [Container Registry ](https://console.cloud.google.com/gcr/images) 
 * bauen Sie das Docker Image mit `make docker-rrimage`
 * verschaffen Sie sich Zugriff auf ihren Cluster mit  
   ` gcloud container clusters get-credentials [CLUSTERNAME] --zone [ZONE] --project [PROJECT_ID]`
 * Das Ausrollen erfolgt mit  
   `kubectl apply -f todo-kubernetes.yaml`
 * Die IP Adresse des Loadbalancers können Sie mit `kubectl get ing` ermitteln.
 * Testen Sie den HTTP-Zugriff

### Cloud Functions
  Wir erstellen zwei Cloud Functions, eine für das automatische Erstellen von Thumbnails in Cloud Storage und eine HTTP Funktion zum umwandeln von Markdown in HTML
  * [Erstellen Sie einen Bucket in Cloud Storage](https://console.cloud.google.com/storage/create-bucket)
  * im Ordner cloudfunctions/imagemagic können Sie die Funktion mit  
    `gcloud functions deploy resize_image --trigger-bucket=gs://$TEST_BUCKET/ --runtime python37` installieren 
  * laden Sie ein beliebiges Bild hoch und testen Sie das Ergebnis im Cloud Storage Browser
* Für die `markdown_to_html` Demo benötigen Sie das [pandoc binary](https://github.com/jgm/pandoc/releases/download/2.5/pandoc-2.5-linux.tar.gz)
* wenn sie `pandoc` extrahiert haben können Sie die Funktion mit `gcloud functions deploy markdown_to_html --runtime python37 --trigger-http` installieren
* in der Ausgabe finden Sie nun die URL ihrer Funktion
* Testen Sie die Konvertierung mit:  
  `curl -sF'doc=@test.md' [URL] `
  
### Cloud Build 
  Wir bauen todo in Cloud Build und deployen in Kubernetes
 * Repository anlegen auf der Kommandozeile: `gcloud source repos create todo` 
 * fügen Sie das gcloud Repo als Remote hinzu: `git remote add gcloud <URL>` (Die Rek-Url können Sie mit `gcloud source repos list` ermitteln)
 * `git push gcloud`
 * Erstellen Sie einen [Build Trigger](https://console.cloud.google.com/cloud-build/triggers/add) für das erstellte Repostitory 
 * für das Erstellen des Docker-Images genügt es wenn Sie in der Buildkonfiguration `Dockerfile` auswählen
 * Testen Sie anschließend die Buildkonfiguration `cloudbuild.yaml`, damit könnes Sie Build Steps festlegen und das erstellte Dockerimage in Kubernetes bereitstellen.
 

  

### todo build 

```
go build

./todo

#get todos
curl 127.0.0.1:8080

#add todo
curl 127.0.0.1:8080 -dtodo="prepare for workshop 💥"

#with postgres db:
docker run --rm -it -p 5432:5432 -e POSTGRES_PASSWORD=geheim postgres
PGPASSWORD=geheim psql -h127.0.0.1 -Upostgres postgres -ac 'CREATE DATABASE test;'

./todo -addr 0.0.0.0:8080 -db "postgres://postgres:geheim@localhost/test?sslmode=disable"

```
