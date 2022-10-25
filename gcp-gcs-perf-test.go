package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

const (
	appName    = "gcp-gcs-perf-stat"
	appVersion = "0.1.3"
)

var (
	cfg = config{}

	localTime time.Time
	tzOffset  int

	logger = logrus.New()
	log    = logger.WithFields(logrus.Fields{})
)

func main() {
	// chShutdown := make(chan os.Signal, 1)

	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Creates a Bucket instance.
	// bucket := client.Bucket(bucketName)

	w := log.WithFields(logrus.Fields{"bucket": cfg.BucketName}).Writer()
	defer w.Close()

	_, err = getBucketMetadata(w, cfg.BucketName)
	if err != nil {
		log.WithFields(logrus.Fields{"err": err}).Fatal("failed to get bucket metadata")
	}

	// wait group
	var wg sync.WaitGroup

	for i := 0; i < cfg.NumObjects; i++ {
		log.WithFields(logrus.Fields{"thread": i}).Info("starting thread")
		wg.Add(1)
		go worker(w, i, &wg)
	}

	log.Info("waiting for workers to finish")
	wg.Wait()

	log.WithFields(logrus.Fields{"timeElapsed": time.Since(localTime).Seconds()}).Info("time elapsed since start (seconds)")
	logger.Exit(0)
}

// worker
func worker(w io.Writer, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	objectName := uuid.Must(uuid.NewV4()).String()

	log.WithFields(logrus.Fields{"thread": id, "objectName": objectName}).Info("thread started")

	err := streamFileUpload(w, cfg.BucketName, objectName, "instance: "+cfg.InstanceGUID+" thread: "+strconv.Itoa(id))
	if err != nil {
		log.WithFields(logrus.Fields{"err": err}).Fatal("failed to write to object")
	}
	log.WithFields(logrus.Fields{"thread": id, "objectName": objectName}).Info("thread finished")
}

// streamFileUpload uploads an object via a stream.
func streamFileUpload(w io.Writer, bucket, object, data string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	b := []byte(data)
	buf := bytes.NewBuffer(b)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.

	if _, err = io.Copy(wc, buf); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	// Data can continue to be added to the file until the writer is closed.
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	// fmt.Fprintf(w, "%v uploaded to %v.\n", object, bucket)

	return nil
}

// getBucketMetadata gets the bucket metadata.
func getBucketMetadata(w io.Writer, bucketName string) (*storage.BucketAttrs, error) {
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, cfg.ConnectionTimeout)
	defer cancel()
	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	fmt.Fprintf(w, "BucketName: %v\n", attrs.Name)
	fmt.Fprintf(w, "Location: %v\n", attrs.Location)
	fmt.Fprintf(w, "LocationType: %v\n", attrs.LocationType)
	fmt.Fprintf(w, "StorageClass: %v\n", attrs.StorageClass)
	fmt.Fprintf(w, "Turbo replication (RPO): %v\n", attrs.RPO)
	fmt.Fprintf(w, "TimeCreated: %v\n", attrs.Created)
	fmt.Fprintf(w, "Metageneration: %v\n", attrs.MetaGeneration)
	fmt.Fprintf(w, "PredefinedACL: %v\n", attrs.PredefinedACL)
	if attrs.Encryption != nil {
		fmt.Fprintf(w, "DefaultKmsKeyName: %v\n", attrs.Encryption.DefaultKMSKeyName)
	}
	if attrs.Website != nil {
		fmt.Fprintf(w, "IndexPage: %v\n", attrs.Website.MainPageSuffix)
		fmt.Fprintf(w, "NotFoundPage: %v\n", attrs.Website.NotFoundPage)
	}
	fmt.Fprintf(w, "DefaultEventBasedHold: %v\n", attrs.DefaultEventBasedHold)
	if attrs.RetentionPolicy != nil {
		fmt.Fprintf(w, "RetentionEffectiveTime: %v\n", attrs.RetentionPolicy.EffectiveTime)
		fmt.Fprintf(w, "RetentionPeriod: %v\n", attrs.RetentionPolicy.RetentionPeriod)
		fmt.Fprintf(w, "RetentionPolicyIsLocked: %v\n", attrs.RetentionPolicy.IsLocked)
	}
	fmt.Fprintf(w, "RequesterPays: %v\n", attrs.RequesterPays)
	fmt.Fprintf(w, "VersioningEnabled: %v\n", attrs.VersioningEnabled)
	if attrs.Logging != nil {
		fmt.Fprintf(w, "LogBucket: %v\n", attrs.Logging.LogBucket)
		fmt.Fprintf(w, "LogObjectPrefix: %v\n", attrs.Logging.LogObjectPrefix)
	}
	if attrs.CORS != nil {
		fmt.Fprintln(w, "CORS:\n")
		for _, v := range attrs.CORS {
			fmt.Fprintf(w, "\tMaxAge: %v\n", v.MaxAge)
			fmt.Fprintf(w, "\tMethods: %v\n", v.Methods)
			fmt.Fprintf(w, "\tOrigins: %v\n", v.Origins)
			fmt.Fprintf(w, "\tResponseHeaders: %v\n", v.ResponseHeaders)
		}
	}
	if attrs.Labels != nil {
		fmt.Fprintf(w, "Labels:\n")
		for key, value := range attrs.Labels {
			fmt.Fprintf(w, "\t%v = %v\n", key, value)
		}
	}
	return attrs, nil
}
