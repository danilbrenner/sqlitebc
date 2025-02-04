package storage

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"log"
	"os"
)

type azureBackupStorage struct {
	accountName   string
	accountKey    string
	containerName string
}

func NewAzureBackupStorage(accountName string, accountKey string, containerName string) *azureBackupStorage {
	return &azureBackupStorage{accountName: accountName, accountKey: accountKey, containerName: containerName}
}

func getClient(s *azureBackupStorage) (*azblob.Client, error) {
	cred, err := azblob.NewSharedKeyCredential(s.accountName, s.accountKey)
	if err != nil {
		return nil, err
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", s.accountName)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (s *azureBackupStorage) Upload(directory string, backupFileName string) error {
	blobName := fmt.Sprintf("%s/%s", directory, backupFileName)
	client, err := getClient(s)
	if err != nil {
		return err
	}

	file, err := os.Open(backupFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = client.UploadFile(context.TODO(), s.containerName, blobName, file, nil)
	if err != nil {
		return err
	}

	log.Printf("Uploaded %s to %s", backupFileName, s.containerName)

	return nil
}

func getBlobs(s *azureBackupStorage, client *azblob.Client, directory string) (*[]backupItem, error) {
	var allItems []backupItem

	pager := client.NewListBlobsFlatPager(s.containerName, &azblob.ListBlobsFlatOptions{Prefix: &directory})
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		for _, blob := range page.Segment.BlobItems {
			if blob.Deleted == nil || !*blob.Deleted {
				allItems = append(allItems, backupItem{name: *blob.Name, modified: *blob.Properties.LastModified})
			}
		}
	}

	return &allItems, nil
}

func (s *azureBackupStorage) AdjustCapacity(directory string, capacity int) error {
	client, err := getClient(s)
	if err != nil {
		return err
	}

	blobs, err := getBlobs(s, client, directory)
	if err != nil {
		return err
	}

	blobsToRemove := sortOutItemsToRemove(capacity, blobs, []backupItem{})
	for _, blob := range blobsToRemove {
		_, err = client.DeleteBlob(context.TODO(), s.containerName, blob.name, nil)
		if err != nil {
			return err
		}
		log.Printf("Blob removed %v", blob)
	}

	return nil
}
