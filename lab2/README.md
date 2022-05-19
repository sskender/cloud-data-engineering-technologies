# Lab 2 additional notes

## Azure Key Vault Go

### Login to Azure

```bash
az login
```

### Create the Go package

```bash
mkdir go-vault
cd go-vault
go mod init go-vault
go get "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
go get "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
```

### Run Go program

Key Vault name: `kv-tim2-azv-ferlab` \
Secret name: `sa-key1`

```bash
export KEY_VAULT_NAME=kv-tim2-azv-ferlab
export SECRET_NAME=sa-key1
go run main.go
```

## Resources

- [Quickstart: Manage secrets by using the Azure Key Vault Go client library](https://docs.microsoft.com/en-us/azure/key-vault/secrets/quick-create-go)
- [Authentication with the Azure SDK for Go using a managed identity](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication-managed-identity?tabs=azure-cli)
- [Interacting with Azure Key Vault in Go](https://moiaune.dev/2021/08/20/interacting-with-azure-key-vault-in-go/)
