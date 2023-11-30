package dao

import (
	"context"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/objectstorage"
	"io"
	"log"
)

const TenancyOCID string = "your-tenancy-ocid"
const UserOCID string = "your-user-ocid"
const Region string = "your-region"
const Fingerprint string = "your-fingerprint"
const compartmentOCID string = "your-compartment-ocid"
const PrivateKey string = "your-private-key"
const Namespace string = "your-namespace"
const BucketName string = "your-bucket-name"
const Url string = "your-bucket-url"

func GetObjectStorageClient() objectstorage.ObjectStorageClient {
	configurationProvider := common.NewRawConfigurationProvider(TenancyOCID, UserOCID, Region, Fingerprint, PrivateKey, nil)
	objectStorageClient, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(configurationProvider)
	fatalIfError(err)
	return objectStorageClient
}

func PutObject(ctx context.Context, c objectstorage.ObjectStorageClient, namespace, bucketname, objectname string, contentLen int64, contentType string, content io.ReadCloser, metadata map[string]string) error {
	request := objectstorage.PutObjectRequest{
		NamespaceName: &namespace,
		BucketName:    &bucketname,
		ObjectName:    &objectname,
		ContentLength: &contentLen,
		ContentType:   &contentType, //"video/mp4" "image/jpeg"
		PutObjectBody: content,
		OpcMeta:       metadata,
	}
	_, err := c.PutObject(ctx, request)
	log.Println("put object: %s", objectname)
	return err
}

func fatalIfError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
