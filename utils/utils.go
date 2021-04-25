package utils

import (
	"context"	
	"fmt"
	"log"
	"os"
	"io"
	"errors"
	"cloud.google.com/go/storage"
    "google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func CheckIfBucketExists(ctx context.Context, client *storage.Client, bucketName string) {
	bucket := client.Bucket(bucketName)
	_,err := bucket.Attrs(ctx)
	if err != nil {
    log.Fatalf("ERROR: %v",err)
	} else {
		log.Println("Bucket exists. OK.")
	}
}

func CheckIfItemExists(jsonPath, projectID, bucketName, fileName string) (bool, error) {
	list, err := ListObjectsInBucket(jsonPath, projectID, bucketName, "silent")
	if err != nil {
		return false, err
	}
	var found bool 
	for _, value := range list {
		if (fileName == value) {
			found = true
		}
	}
	return found, nil
}


func ListObjectsInBucket(jsonPath, projectID, bucketName, flag string) ([]string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
			return nil, err
	}
	CheckIfBucketExists(ctx, client, bucketName)
	//fmt.Printf("%T    %T\n", ctx, client)

	bucket := client.Bucket(bucketName)
	it := bucket.Objects(ctx, nil)
	
	var i int
	var itemList []string
	if (flag != "silent") {
	fmt.Printf("\n Bucket: %s\n =======================\n", bucketName)
	}
	for {
		i +=1 
		attrs, err := it.Next()
		if err == iterator.Done {
				break
		}
		if err != nil {
				return itemList, fmt.Errorf("Bucket(%q).Objects: %v", bucketName, err)
		}
		itemList = append(itemList, attrs.Name)
		if (flag != "silent") {
			fmt.Printf("	Img_%v: 	%v\n", i, attrs.Name)
		}	
	}
	return itemList, nil
}

func UploadFile(jsonPath, bucketName, filePath, uploadedFilename string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
			log.Fatal(err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("problem opening file for gcs")
	}
	defer file.Close()

	sw := client.Bucket(bucketName).Object(uploadedFilename).NewWriter(ctx)

	if _, err := io.Copy(sw, file); err != nil {
		return err
	}

	if err := sw.Close(); err != nil {
		return err
	}

	return nil 
}
