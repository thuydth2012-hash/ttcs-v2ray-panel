package xray

import (
	"fmt"
	"os"
	"os/exec"

	"ttcs/database"
)

const ConfigPath = "/tmp/ttcs-xray.json"

func WriteConfig() error {
	var inbounds []database.Inbound

	if err := database.DB.Find(&inbounds).Error; err != nil {
		return fmt.Errorf(
			"failed to load inbounds: %w",
		)
	}

	config, err := BuildConfig(inbounds)
	if err != nil {
		return fmt.Errorf(
			"failed to build Xray config: %w",
			err,
		)
	}

	tempPath := ConfigPath + ".tmp"

	if err := os.WriteFile(
		tempPath,
		config,
		0644,
	); err != nil {
		return fmt.Errorf(
			"failed to write temporary config: %w",
			err,
		)
	}

	if err := os.Rename(
		tempPath,
		ConfigPath,
	); err != nil {
		os.Remove(tempPath)

		return fmt.Errorf(
			"failed to replace Xray config: %w",
			err,
		)
	}

	return nil
}
func Apply() error {
	if err := WriteConfig(); err != nil {
		return err
	}

	cmd := exec.Command(
		"sudo",
		"/usr/local/bin/ttcs-xray-apply",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to apply Xray config: %w: %s",
			err,
			string(output),
		)
	}

	return nil
}
