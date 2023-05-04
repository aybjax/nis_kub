# Go microservices

### Go / GoKit / Gorilla mux Http / Protobuffer / Grpc / Redis / Nats Jetstream / Postgres / MongoDB / Kong / Kubernetes / Prometheus / Grafana

## File/Directory/Project Contents
- Part of functions, types and interfaces that are common between services are located [here](https://github.com/aybjax/nis_lib)

- ### /courses
  - _/app/adapter_ contains: environmental variables, "adapters"/"drivers" for queue and cache services (_nats and redis both can be used as either queue or key-value store_). Caching is performed by **Redis** because it is in-memory (instant). However, caching stores for 5 minutes, as caches can clog RAM. Queuing is performed by **Nats Jetstream** as it is small and does not have overhead: no Zookeeper or JVM as in Kafka.
  - _/app_ contains: higher level types for queue and cache: they facade to application interfaces. It also contains interface for service interface.
  - _/app_db_: database connection and database models. Database models has db-specific functionalities, validation and can be converted to and from protobuf DTOs.
    In this service, **mongodb** is used as database of choice. 1. It is scalable. Number of courses can grown unpredicatbly. 2. Each different courses (in class or elective or externships) can have different structure and during the development and production its structure may change further. 3. Speed. NoSql databases are usually faster
  - _/service_: contain implementation of app service. It contains (according to hexagonal architecture, aka Ports and Adapters (https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))):
    - **Application core**: CourseService interface implementations (_ and, in my opinion, middlewares for logging, caching and metrics). **Logging middleware** logs timing, type (grpc/http) and duration of service methods. It writes to Stdout file, for it to be collected by **Fluent bit** and sent to **Loki**. Stdout is used because containers allow easy access to stdout logs. **Caching** middleware stores and servers GET requests and invalidates the cache on mutation requests. **Metrics middleware** collects count (request count) and summary (request latency) metrics. It is then served to prometheus by /metrics http endpoint.
    - **Adapters**: mapping_*.go and endpoints.go files. They convert input arguments from any port, _be it GRPC or HTTP server_, types service types. In our case it is protobuffer DTOs
    - **Ports**: communication with outside world. Here are only GRPC and HTTP handler files. Unfortunately, NATS jetstream transport is not supported by GoKit library (although nats is supported)
      - HTTP server: for external communication thru API gateway
      - GRPC server: for inter-service communication as it is very fast, but requires structured objects
      - Nats jetstream: non-instant inter-service communication. GRPC is used with GET requests, when result is needed instantly. **Queue groups** in Nats Jetstream allows one operation at a time, so that we can remove load on services
  - _contructors.go_: factory methods for dependency injection
  - Dockerfile for containerization
  
- ### /students
  - It displays same structure
  - **Postgres** is used as database of choice. 1. Postgres is fast 2. Students number is usually predictable and their information has known structure. 3. Sql databases allow joins between multiple tables. And students table can be used further with other services (_other services also store in the same db_). For example it can be joined with roles table to used in authentication services. Furthermore, it can also be joined with organizations table to give functionality related to partitioning by organization
  - Dockerfile for containerization
  
- ### /initializer
  - It's only function is to connect to NATS Jetstream and manage streams. Streams can be created only once, and connection are kept until program is running. So, if student or course pods manage streams, it would be impossible to operate and one of the duplicated streams would create streams and others would panic. Also it would be impossible to control who is in charge of streams
    - It makes students and courses services scalable
  - Dockerfile for containerization
  
- ### /kong
  - Just used in developmental environment for docker-compose configurations
- ### /k8s
  - configuration and helm value yaml files for kubernetes. Docker in docker minikube node is used as kubernetes
    - _/courses_, _/initializer_, _/mongo_, _/nats_, _/postgres_, _/redis_, _/students_ contain **configMaps**, **secretMaps**, **deployments** and **services**.
    - _/metrics_ directory contains kong plugins and values files:
      - **kong.yml** contains annotations for kong helm-chart for prometheus to collect metrics
      - **auth.plugin.yml** contains key-auth plugin and consumer for kong. Key-auth is used for its simplicity.
      - **password.yml** contains secretMap for kong consumer. Head or body must contain x-api-key or apikey key with _T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V_ value.
      - **prometheus-kong-plugin.yml** contains plugin for prometheus exporter
      - **prometheus.yml** helm-chart values for prometheus
      - **grafana.yml** contains helm-chart values for grafana. It is only used for testing. During gathering of chart, I used grafana web app with port-forwarded _(from kubernetes to localhost)_ and tunneled _(from localhost to www)_ url.
      - **ingress.yml** contains configurations for kong ingress controller. It allows communication between Kong API gateway and kubernetes clusters
      
## How it works
  - Each service exposes simple CRUD apis. Middleware caches GET requests and serves it in subsequent requests. However, mutating endpoints invalidate the cache. Get requests about another service _(get student courses as example)_ implemented with the use of GRPC servers. However, Mutating requests are conveyed thru queues, because it requires data transformation, validation and cleanup if request was invalid. Also it allows 1 operation at a time, so that it does not load the service, although data may take time to be written

## Results
  - Docker images are located at **/courses/Dockerfile**, **/students/Dockerfile** and **/initializer/Dockerfile**. Images are also located at dockerhub repository of user aybjax.
    - There is docker-compose.yml file for fast developmental experience
  - Kubernetes config files and helm-charts values files are located in **/k8s** directory.
  - Initiating project in minikube required following steps, assuming you are in **/k8s** directory:
    1. ```minikube start```
    2. ```helm repo add kong https://charts.konghq.com```
    3. ```helm repo add prometheus-community https://prometheus-community.github.io/helm-charts```
    4. ```helm repo add grafana https://grafana.github.io/helm-charts```
    5. ```helm repo update```
    6. ```helm upgrade --install prometheus prometheus-community/prometheus --values ./metrics/prometheus.yml```
    7. ```helm upgrade --install grafana grafana/grafana --values ./metrics/grafana.yml```
    8. ```helm upgrade --install loki-stack grafana/loki-stack --set fluent-bit.enabled=true,promtail.enabled=false```
    9. ```helm upgrade --install kong kong/kong --values ./metrics/kong.yml --set admin.enabled=true --set=admin.http.enabled=true --set ingressController.installCRDs=false```
    10. ```kubectl apply -f ./metrics/prometheus-kong-plugin.yml```
    11. ```kubectl apply -f ./metrics/ingress.yml```
    12. ```kubectl apply -f ./metrics/password.yml```
    13. ```kubectl apply -f ./metrics/auth.plugin.yml```
    14. ```kubectl apply -f ./mongo/mongo-config.yml```
    15. ```kubectl apply -f ./mongo/mongo-secret.yml```
    16. ```kubectl apply -f ./mongo/mongo.yml```
    17. ```kubectl apply -f ./postgres/postgres-config.yml```
    18. ```kubectl apply -f ./postgres/postgres-secret.yml```
    19. ```kubectl apply -f ./postgres/postgres.yml```
    20. ```kubectl apply -f ./redis/redis-config.yml```
    21. ```kubectl apply -f ./redis/redis.yml```
    22. ```kubectl apply -f ./nats/nats-config.yml```
    23. ```kubectl apply -f ./nats/nats.yml```
    24. ```kubectl apply -f ./initializer/initializer.yml```
    25. ```kubectl apply -f ./courses/courses-config.yml```
    26. ```kubectl apply -f ./courses/courses.yml```
    27. ```kubectl apply -f ./students/students-config.yml```
    28. ```kubectl apply -f ./students/students.yml```
    29. ```kubectl port-forward deployment.apps/prometheus-server 9090```
    30. ```kubectl port-forward deployment.apps/grafana 3000```
    31. ```kubectl port-forward deployment.apps/kong-kong 8000```
    32. ```echo -e "Grafana password is:\n$GRAFANA_PASS\n"```: username is _admin_
  - You can also use minikube tunnels for API Gateways
  - **NB** you can also apply directories, not each file
- Prometheus/Loki/Fluent bit/Grafana results:
