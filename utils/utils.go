package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func CheckIfBucketExists(ctx context.Context, client *storage.Client, bucketName string) {
	bucket := client.Bucket(bucketName)
	_, err := bucket.Attrs(ctx)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
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
	for index, _ := range list {
		if fileName == index {
			found = true
		}
	}
	return found, nil
}

func ListObjectsInBucket(jsonPath, projectID, bucketName, flag string) (map[string]time.Time, error) {
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
	itemList := make(map[string]time.Time)
	if flag != "silent" {
		fmt.Printf("\n Bucket: %s\n =======================\n", bucketName)
	}
	for {
		i += 1
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return itemList, fmt.Errorf("Bucket(%q).Objects: %v", bucketName, err)
		}
		itemList[attrs.Name] = attrs.Created
		if flag != "silent" {
			fmt.Printf("	Img_%v: 	%v    Created_at: %v\n", i, attrs.Name, attrs.Created)
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

func CheckifImgUpdated(imglist map[string]time.Time, downloadSuseLink string) ([]string, error) {
	err := CheckNetworkFine(downloadSuseLink)
	if err != nil {
		return nil, err
	}
	cmd := []string{"curl", downloadSuseLink}
	output, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return nil, err
	}
	stroutput := string(fmt.Sprintf("%s", output))
	var imgToUpdate []string
	if len(imglist) > 0 {
		for index, value := range imglist {
			//fmt.Println(value)
			if strings.Contains(stroutput, index) {
				log.Printf("%s is the latest Google Cloud Image.\n", index)
			} else {
				day, _ := ImgVersioningParser(stroutput, index, "x86_64")
				if day.Sub(value).Hours() > 3 {
					imgToUpdate = append(imgToUpdate, index)
				}
			}
		}
	}
	return imgToUpdate, nil
}

func ImgVersioningParser(webpage string, image string, arch string) (time.Time, error) {
	tmpwebpageslice := strings.Split(webpage, "Details")
	imgprefix := strings.Split(image, fmt.Sprintf(".%s-", arch))[0]
	imgregister := make(map[string][]string)
	regserver := regexp.MustCompile(".tar.gz\"")
	var day time.Time
	var err error
	for _, value := range tmpwebpageslice {
		if regserver.FindString(value) != "" && strings.Contains(value, arch) && strings.Contains(value, "GCE") {
			if strings.Contains(value, imgprefix) {
				day, err = ParseWebHTMLLine(value)
				if err != nil {
					return day, err
				}
				imgregister[imgprefix] = append(imgregister[imgprefix], strings.Split(value, fmt.Sprintf(".%s-", arch))[1])
			}
		}
	}
	//fmt.Println(imgregister)
	return day, nil
}

func ParseWebHTMLLine(htmlLine string) (time.Time, error) {
	fmt.Println(htmlLine)
	reg := regexp.MustCompile(`\d{2}-\w{3,9}-\d{4} \d{2}:\d{2}`)
	timestamp := fmt.Sprintf("%s CET", reg.FindStringSubmatch(htmlLine)[0])
	if strings.Contains(timestamp, "-202") {
		timestamp = strings.Replace(timestamp, "20", "", 1)
	}
	day, err := time.Parse(time.RFC822, strings.Replace(timestamp, "-", " ", 10))
	if err != nil {
		return day, err
	}
	reg = regexp.MustCompile(`\d{4}-\d{7}`)
	return day, nil
}

func ReplaceImagesOnGCE(imgToUpdate []string, jsonPath, bucketName, downloadSuseLink string) error {

	return nil
}

func CheckNetworkFine(downloadSuseLink string) error {
	url, err := url.Parse(downloadSuseLink)
	if err != nil {
		return err
	}
	//fmt.Println(url.Host)
	out, err := exec.Command("which", "ping").CombinedOutput()
	if err != nil {
		return err
	}
	if strings.Contains(fmt.Sprintf("%s", string(out)), "bin/ping") {
		out, err := exec.Command("ping", "-c", "1", "-W", "1", url.Host).CombinedOutput()
		if err != nil {
			if strings.Contains(fmt.Sprintf("%s", string(out)), "100% packet loss") {
				err = errors.New("Error: " + fmt.Sprintf("%s", string(out)))
			}
			return err
		}
	} else {
		out, err := exec.Command("which", "fping").CombinedOutput()
		if err != nil {
			return err
		}
		if strings.Contains(fmt.Sprintf("%s", string(out)), "bin/fping") {
			out, err := exec.Command("fping", "-c1", "-t500", url.Host).CombinedOutput()
			if err != nil {
				if strings.Contains(fmt.Sprintf("%s", string(out)), "100% packet loss") {
					err = errors.New("Network Error: " + fmt.Sprintf("%s", string(out)))
				}
				return err
			}
		}
	}
	return err
}
