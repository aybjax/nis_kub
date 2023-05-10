minikube start

helm repo add kong https://charts.konghq.com
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts

helm repo update

# Helm 3
helm upgrade --install prometheus prometheus-community/prometheus --values ./metrics/prometheus.yml
helm upgrade --install grafana grafana/grafana --values ./metrics/grafana.yml
# helm upgrade --install fluent-bit grafana/fluent-bit \
#     --set loki.serviceName=loki.svc.cluster.local # birge ili net?
helm upgrade --install loki-stack grafana/loki-stack \
    --set fluent-bit.enabled=true,promtail.enabled=false
helm upgrade --install kong kong/kong --values ./metrics/kong.yml --set admin.enabled=true --set=admin.http.enabled=true --set ingressController.installCRDs=false

kubectl apply -f ./metrics/prometheus-kong-plugin.yml

kubectl apply -f ./metrics/ingress.yml
kubectl apply -f ./metrics/password.yml
kubectl apply -f ./metrics/auth.plugin.yml

kubectl apply -f ./mongo/mongo-config.yml
kubectl apply -f ./mongo/mongo-secret.yml
kubectl apply -f ./mongo/mongo.yml

kubectl apply -f ./postgres/postgres-config.yml
kubectl apply -f ./postgres/postgres-secret.yml
kubectl apply -f ./postgres/postgres.yml

kubectl apply -f ./redis/redis-config.yml
kubectl apply -f ./redis/redis.yml

kubectl apply -f ./nats/nats-config.yml
kubectl apply -f ./nats/nats.yml

kubectl apply -f ./initializer/initializer.yml

kubectl apply -f ./courses/courses-config.yml
kubectl apply -f ./courses/courses.yml

kubectl apply -f ./students/students-config.yml
kubectl apply -f ./students/students.yml

# # does not work on windows => use kubectl get pods and kubectl port-forward NAME PORT manually
# # POD_NAME=$(kubectl get pods -l "app=prometheus,component=server" -o jsonpath="{.items[0].metadata.name}")
# # kubectl port-forward $POD_NAME 9090
# kubectl port-forward deployment.apps/prometheus-server 9090
# # POD_NAME=$(kubectl get pods -l "app.kubernetes.io/instance=grafana" -o jsonpath="{.items[0].metadata.name}")
# # kubectl port-forward $POD_NAME 3000
# kubectl port-forward deployment.apps/grafana 3000
# # POD_NAME=$(kubectl get pods -o jsonpath="{.items[0].metadata.name}")
# # kubectl port-forward $POD_NAME 8000
# kubectl.exe port-forward deployment.apps/kong-kong 8000

# GRAFANA_PASS=$(kubectl get secret grafana -o jsonpath="{.data.admin-password}" | base64 --decode );
# #UojzGRAE8gVVK5HayMjAKx6I1jplw3A853qc3HOQ

# echo -e "Grafana password is:\n$GRAFANA_PASS\n"

# minikube tunnel