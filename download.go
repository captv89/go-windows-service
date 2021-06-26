package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-file-go/azfile"
)

// Please set environment variable ACCOUNT_NAME and ACCOUNT_KEY to your storage accout name and account key,
// before run the examples.
func accountInfo() (string, string) {
	return os.Getenv("ACCOUNT_NAME"), os.Getenv("ACCOUNT_KEY")
}

func download(foldername string, filename string) (e error) {
	// From the Azure portal, get your Storage account file service URL endpoint.
	accountName, accountKey := accountInfo()
	fmt.Println(accountKey, accountName)
	credential, err := azfile.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}

	// All HTTP operations allow you to specify a Go context.Context object to control cancellation/timeout.
	ctx := context.Background() // This example uses a never-expiring context.

	// From the Azure portal, get your Storage account file service URL endpoint.
	// The URL typically looks like this:
	u, _ := url.Parse(fmt.Sprintf("https://%s.file.core.windows.net/%s/%s", accountName, foldername, filename))

	// Create a URL that references a to-be-created file in your Azure Storage account's directory.
	// This returns a FileURL object that wraps the file's URL and a request pipeline (inherited from directoryURL)
	fileURL := azfile.NewFileURL(*u, azfile.NewPipeline(credential, azfile.PipelineOptions{})) // File names can be mixed case and is case insensitive
	fmt.Println("Trying", fileURL)

	// Query the file's properties and metadata
	get, err := fileURL.GetProperties(ctx)
	if err != nil {
		fmt.Println(err)
		e = err
		return e
	}

	if get.StatusCode() == 200 {

		// // Show some of the file's read-only properties
		// fmt.Println(get.FileType(), get.ETag(), get.LastModified())

		// Trigger download.
		downloadResponse, err := fileURL.Download(ctx, 0, azfile.CountToEnd, false) // 0 offset and azfile.CountToEnd(-1) count means download entire file.
		if err != nil {
			log.Fatal(err)
		}

		contentLength := downloadResponse.ContentLength() // Used for progress reporting to report the total number of bytes being downloaded.
		fmt.Println(contentLength)

		// Setup RetryReader options for stream reading retry.
		retryReader := downloadResponse.Body(azfile.RetryReaderOptions{MaxRetryRequests: 3})

		// NewResponseBodyStream wraps the RetryReader with progress reporting; it returns an io.ReadCloser.
		progressReader := pipeline.NewResponseBodyProgress(retryReader,
			func(bytesTransferred int64) {
				fmt.Printf("Downloaded %d of %d bytes.\n", bytesTransferred, contentLength)
			})
		defer progressReader.Close() // The client must close the response body when finished with it

		file, err := os.Create(filename) // Create the file to hold the downloaded file contents.
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		written, err := io.Copy(file, progressReader) // Write to the file by reading from the file (with intelligent retries).
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Download Completed!")
		}
		_ = written // Avoid compiler's "declared and not used" error
	}
	return
}
