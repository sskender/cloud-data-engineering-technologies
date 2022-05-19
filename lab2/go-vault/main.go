package main

import (
    "context"
    "fmt"
    "log"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

func vaultInfo() {
	keyVaultName := "kv-tim2-azv-ferlab"
	keyVaultUrl := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)
	secretName := "sa-key1"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	client, err := azsecrets.NewClient(keyVaultUrl, cred, nil)
    if err != nil {
        log.Fatalf("failed to create a client: %v", err)
    }

	getResp, err := client.GetSecret(context.TODO(), secretName, nil)
	if err != nil {
	  log.Fatalf("failed to get the secret: %v", err)
	}

	fmt.Printf("secretValue: %s\n", *getResp.Value)
}

func main() {
    vaultInfo()
}
