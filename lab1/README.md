# Lab 1 additional notes

## Homebrew installs

```bash
brew install go && \
brew tap azure/functions && \
brew install azure-functions-core-tools@4 && \
brew install azure-cli
```

## Azure functions

### Python

Open folder [`tim2-af-hello-world-py`](./tim2-af-hello-world-py) in VS Code.

Run function locally:

```bash
func start
```

### GoLang

Open folder [`tim2-af-hello-world-go`](./tim2-af-hello-world-go)  in VS Code.

Run function locally:

```bash
go build helloworld.go
func start
```

Build and run Docker locally:

```bash
docker build --platform linux/amd64 --tag tim2-af-hello-world-go .  # when building on M1 Mac
docker run --platform linux/amd64 -p 8080:8080 tim2-af-hello-world-go  # when building on M1 Mac
curl "http://localhost:8080/api/tim2-af-hello-world-go"
```

## AKS deployment

View [docs](./docs/) for setting up AKS and ACR.

Push Docker image to ACR:

```bash
az login
az acr login --name tim2acrferlab
docker tag tim2-af-hello-world-go:latest tim2acrferlab.azurecr.io/tim2-af-hello-world-go:latest
docker push tim2acrferlab.azurecr.io/tim2-af-hello-world-go:latest
```

Connect to AKS:

```bash
az aks get-credentials --resource-group tim2-rg-ferlab --name tim2-aks-ferlab
kubectl get nodes
```

Deploy to AKS:

```bash
kubectl apply -f ./manifest.yaml
kubectl apply -f ./service.yaml
kubectl get services
```

## Resources

- [Quickstart: Create a Go or Rust function in Azure using Visual Studio Code](https://docs.microsoft.com/en-us/azure/azure-functions/create-first-function-vs-code-other?tabs=go%2Cmacos)
- [Quickstart: Deploy an Azure Kubernetes Service cluster using the Azure CLI](https://docs.microsoft.com/en-us/azure/aks/learn/quick-kubernetes-deploy-cli)
- [Deploy a containerized application on Azure Kubernetes Service](https://docs.microsoft.com/en-us/learn/modules/aks-deploy-container-app/)
- [Azure — Deploying Vue App With NGINX on AKS](https://medium.com/bb-tutorials-and-thoughts/azure-deploying-vue-app-with-nginx-on-aks-530e974daf1e)

## Results

### tim2-af-hello-world-py Azure Function

[https://tim2-af-hello-world-py.azurewebsites.net/api/tim2-af-hello-world-py](https://tim2-af-hello-world-py.azurewebsites.net/api/tim2-af-hello-world-py)

### tim2-af-hello-world-go AKS

[http://20.223.80.213/api/tim2-af-hello-world-go](http://20.223.80.213/api/tim2-af-hello-world-go)
