package migration

import (
	"archive/tar"
	"bytes"
	"database/sql"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"

	localTypes "github.com/justsushant/envbox/types"
)

func MigrateDockerImgUp(cli *client.Client, db *sql.DB, migTable []localTypes.ImageMigration) error {
	for _, mig := range migTable {
		log.Printf("loading %q image\n", mig.Name)

		// check for image on local
		ok, err := checkIfImageExists(cli, mig.Name)
		if err != nil {
			return fmt.Errorf("error while checking if image exists: %v", err)
		}
		if ok {
			log.Printf("image %q exists\n", mig.Name)

			// if image already exists, check if it is already loaded in db
			err = insertImageRecordInDB(db, mig.Label, mig.Name)
			if err != nil {
				log.Fatalf("error while inserting image record in db: %v", err)
			}
			continue
		}

		// if not found, load image
		resp, err := loadImage(cli, mig.Name, mig.Path)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// show output of build process
		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			log.Fatalf("error while reading build process response: %v", err)
		}

		// inserting image record in db
		err = insertImageRecordInDB(db, mig.Label, mig.Name)
		if err != nil {
			log.Fatalf("error while inserting image record in db: %v", err)
		}

		log.Printf("image %q loaded successfully\n", mig.Name)
	}

	return nil
}

func loadImage(cli *client.Client, dockerFileName, dockerFilePath string) (types.ImageBuildResponse, error) {
	// reading dockerfile
    reader, err := os.Open(dockerFilePath)
    if err != nil {
        return types.ImageBuildResponse{}, fmt.Errorf("error while opening docker file: %v", err)
    }
    dockerFileReader, err := io.ReadAll(reader)
    if err != nil {
        return types.ImageBuildResponse{}, fmt.Errorf("error while reading docker file: %v", err)
    }

	// creating tar reader
	buf := new(bytes.Buffer)
    tw := tar.NewWriter(buf)
    defer tw.Close()

    tarHeader := &tar.Header{
        Name: dockerFileName,
        Size: int64(len(dockerFileReader)),
    }

	// writing to tar writer
    err = tw.WriteHeader(tarHeader)
    if err != nil {
        return types.ImageBuildResponse{}, fmt.Errorf("error while reading tar header: %v", err)
    }
    _, err = tw.Write(dockerFileReader)
    if err != nil {
        return types.ImageBuildResponse{}, fmt.Errorf("error while writing tar file: %v", err)
    }
    dockerFileTarReader := bytes.NewReader(buf.Bytes())

	// building docker image from tar
    resp, err := cli.ImageBuild(
        context.Background(),
        dockerFileTarReader,
        types.ImageBuildOptions{
			Context: dockerFileTarReader,	
			Tags: 	 []string{dockerFileName},
            Dockerfile: dockerFileName,
            Remove:     true,	// containers created during build process will be removed
			ForceRemove: true, // cleanup if build process fails
	})
    if err != nil {
        return types.ImageBuildResponse{}, fmt.Errorf("error while building docker image: %v", err)
    }

	return resp, nil
}

func checkIfImageExists(cli *client.Client, imageName string) (bool, error) {
	ctx := context.Background()

	// listing image name via reference
	img, err := cli.ImageList(ctx, image.ListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", imageName)),
	})
	if err != nil {
		return false, fmt.Errorf("error while listing images: %v", err)
	}

	if len(img) > 0 {
		return true, nil
	}
	return false, nil
}