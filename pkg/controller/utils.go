package controller

import (
	"io/ioutil"
	"os"
)

func writeKubeconfigToTempFile(kubeconfig []byte) (path string, error error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "kubeconfig-")
	if err != nil {
		return "", err
	}

	if _, err = tmpFile.Write(kubeconfig); err != nil {
		return "", err
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}
