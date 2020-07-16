# deploy image service

kubectl delete configmaps image-server-configmap -n cloudai-2
kubectl create configmap image-server-configmap --from-file=server.py --namespace=cloudai-2

kubectl delete -f deploy.yaml
kubectl create -f deploy.yaml

kubectl create -f service.yaml

# deploy image clean

kubectl delete configmaps image-clean-job-configmap -n cloudai-2
kubectl create configmap image-clean-job-configmap --from-file=clean_job.py --namespace=cloudai-2

kubectl delete -f cronjob-clean.yaml
kubectl create -f cronjob-clean.yaml