package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
    "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func vaultInfo() (string) {
	keyVaultName := os.Getenv("KEY_VAULT_NAME")
	keyVaultUrl := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)
	secretName := os.Getenv("SECRET_NAME")

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

	fmt.Printf("Secret value: %s\n", *getResp.Value)

    return *getResp.Value
}

func main() {
    accountName := "satim2ferlab"
    secretKey := vaultInfo()
    containerName := "satim2ferlab"

    cred, err := azblob.NewSharedKeyCredential(accountName, secretKey)
    if err != nil {
        log.Fatalf("failed to get the secret: %v", err)
    }

    url := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
    service, err := azblob.NewServiceClientWithSharedKey(url, cred, nil)
    
    ctx := context.Background()
    container, err := service.NewContainerClient(containerName)

    /*
    _, err = container.Create(ctx, nil)
    if err != nil {
        log.Fatalf("failed to get the secret: %v", err)
    }
    */

    blockBlob, err := container.NewBlockBlobClient("superstore_dataset2011-2015.csv")


    // Download the blob's contents and ensure that the download worked properly
    get, err := blockBlob.Download(ctx, nil)
    fmt.Printf("File: %s\n", get)


    //handle(err)


    //data, err := os.Open("superstore.csv")
    /*
    blobClient, err := container.NewBlockBlobClient("superstore.csv")
    if err != nil {
        log.Fatal(err)
    }

    o := azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: "text/csv",   //  Add any needed headers here
		},
	}
    
    datab := []byte("asd")
    // Upload to data to blob storage
    _, err = blobClient.Upload(ctx, bytes.NewBufferString("asd"), o)
    
    if err != nil {
        log.Fatalf("Failure to upload to blob: %+v", err)
    }

*/
}
