package main

import (
	"bytes"
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gocarina/gocsv"
	"io/ioutil"

	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type Order struct {
	Row_ID         string `csv:"Row ID"`
	Order_ID       string `csv:"Order ID"`
	Order_Date     string `csv:"Order Date"`
	Ship_Date      string `csv:"Ship Date"`
	Ship_Mode      string `csv:"Ship Mode"`
	Customer_ID    string `csv:"Customer ID"`
	Customer_Name  string `csv:"Customer Name"`
	Segment        string `csv:"Segment"`
	City           string `csv:"City"`
	State          string `csv:"State"`
	Country        string `csv:"Country"`
	Postal_Code    string `csv:"Postal Code"`
	Market         string `csv:"Market"`
	Region         string `csv:"Region"`
	Product_ID     string `csv:"Product ID"`
	Category       string `csv:"Category"`
	Sub_Category   string `csv:"Sub-Category"`
	Product_Name   string `csv:"Product Name"`
	Sales          string `csv:"Sales"`
	Quantity       string `csv:"Quantity"`
	Discount       string `csv:"Discount"`
	Profit         string `csv:"Profit"`
	Shipping_Cost  string `csv:"Shipping Cost"`
	Order_Priority string `csv:"Order Priority"`
	NotUsed        string `csv:"-"`
}

func vaultInfo() string {
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

func sendOrderMessageToTopic(order Order, sender *azservicebus.Sender) {

	out, err := json.Marshal(order)
	if err != nil {
		panic(err)
	}

	// send a single message
	err = sender.SendMessage(context.TODO(), &azservicebus.Message{
		Body: []byte(string(out)),
	}, nil)
}

func receiveAllOrders(servicebus_client *azservicebus.Client, topicName string, subscriptionName string) {
	receiver, err := servicebus_client.NewReceiverForSubscription(topicName, subscriptionName, nil)

	if err != nil {
		log.Fatalf("Failed to create the receiver: %s", err.Error())
	}

	// Receive a fixed set of messages. Note that the number of messages
	// to receive and the amount of time to wait are upper bounds.
	messages, err := receiver.ReceiveMessages(context.TODO(),
		// The number of messages to receive. Note this is merely an upper
		// bound. It is possible to get fewer message (or zero), depending
		// on the contents of the remote queue or subscription and network
		// conditions.
		5,
		&azservicebus.ReceiveMessagesOptions{},
	)

	if err != nil {
		log.Fatalf("Failed to get messages: %s", err.Error())
	}

	subscriptionFetchedOrders := []Order{}

	for _, message := range messages {
		data := Order{}
		json.Unmarshal(message.Body, &data)
		subscriptionFetchedOrders = append(subscriptionFetchedOrders, data)

		// For more information about settling messages:
		// https://docs.microsoft.com/azure/service-bus-messaging/message-transfers-locks-settlement#settling-receive-operations
		if err := receiver.CompleteMessage(context.TODO(), message, nil); err != nil {
			log.Printf("Error completing message: %s", err.Error())
		}
	}

	for _, order := range subscriptionFetchedOrders {
		fmt.Printf("\nOrder fetched from subscription: %v", order)
	}
}

func main() {
	accountName := "satim2ferlab"
	secretKey := "apaqpTKszPAkVqdBexPcglm+8GqTdt+iYS0iFWtFYrZ2ovDBI0UPXpXU5wVSUY5E5dWxSTu1k7jGPxwH2FGgBQ==" //vaultInfo()
	containerName := "satim2ferlab"

	cred, err := azblob.NewSharedKeyCredential(accountName, secretKey)
	fmt.Print(secretKey)
	if err != nil {
		log.Fatalf("failed to get the secret: %v", err)
	}

	url := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	service, err := azblob.NewServiceClientWithSharedKey(url, cred, nil)

	ctx := context.Background()
	container, err := service.NewContainerClient(containerName)

	blobClient, err := container.NewBlockBlobClient("superstore.csv")
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile("superstore.csv")
	if err != nil {
		fmt.Print(err)
	}

	o, err := blobClient.UploadBuffer(ctx, b, azblob.UploadOption{})

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("\nUpload return value:   %s\n", *o)
	}

	// Can we use same BlockBlobClient used for upload?
	blockBlob, err := container.NewBlockBlobClient("superstore.csv")

	// // Download the blob's contents and ensure that the download worked properly
	get, err := blockBlob.Download(ctx, nil)
	fmt.Printf("\nFile: %s\n", get)

	// Open a buffer, reader, and then download!
	downloadedData := &bytes.Buffer{}
	// RetryReaderOptions has a lot of in-depth tuning abilities, but for the sake of simplicity, we'll omit those here.
	reader := get.Body(&azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		log.Fatal(err)
	}
	err = reader.Close()
	if err != nil {
		log.Fatal(err)
	}

	orders := []*Order{}

	gocsv.UnmarshalBytes(downloadedData.Bytes(), &orders)

	// example of order
	//fmt.Printf("\nOrders : %+v", orders)

	primaryConnectionString := "Endpoint=sb://tim2-sb-ferlab.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=QEfks5h9uZ3gX5X73+wLUXWuY8EIiFD0uGggDBTTPzA="

	servicebus_client, err := azservicebus.NewClientFromConnectionString(primaryConnectionString, nil)

	if err != nil {
		log.Fatalf("Failed to create Service Bus Client: %s", err.Error())
	}

	topicName := "retaildatatopic"
	subscriptionName := "retailDataSubscription"

	sender, err := servicebus_client.NewSender(topicName, nil)

	if err != nil {
		log.Fatalf("Failed to create Sender: %s", err.Error())
	}

	for i, o := range orders {
		if i == 5 {
			break
		}
		sendOrderMessageToTopic(*o, sender)
	}

	receiveAllOrders(servicebus_client, topicName, subscriptionName)

}
