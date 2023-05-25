package internal

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

type OpsManifestService struct {
	opsManifestRunner OpsManifestRunnerI
	fileIO            FileIO
}

func NewOpsManifestServiceWithRunner(opsManifestRunner OpsManifestRunnerI, fileIO FileIO) (*OpsManifestService, error) {
	return &OpsManifestService{opsManifestRunner: opsManifestRunner, fileIO: fileIO}, nil
}

func (o OpsManifestService) RenderManifest(tileConfig io.Reader, tileMetadata io.Reader) (string, error) {
	f, err := o.fileIO.TempFile("", "metadata")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = o.fileIO.Remove(f.Name())
		_ = f.Close()
	}()

	_, err = io.Copy(f, tileMetadata)
	if err != nil {
		return "", err
	}

	configInput, err := io.ReadAll(tileConfig)
	if err != nil {
		return "", err
	}

	manifest, err := o.opsManifestRunner.GetManifest(string(configInput), f.Name())
	if err != nil {
		return "", fmt.Errorf("Unable to retrieve bosh manifest: %s", err)
	}

	y, err := yaml.Marshal(manifest)
	if err != nil {
		return "", err // un-tested
	}

	return string(y), nil
}
