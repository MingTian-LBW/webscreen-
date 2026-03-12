package android

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"webscreen/utils"
)

func ExecADB(args ...string) error {
	adbPath, err := utils.GetADBPath()
	if err != nil {
		return err
	}
	cmd := exec.Command(adbPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetDevices returns a list of connected devices
func GetDevices() ([]AndroidDevice, error) {
	adbPath, err := utils.GetADBPath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(adbPath, "devices")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var adbDevices []AndroidDevice
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "List of devices attached") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			switch parts[1] {
			case "device":
				adbDevices = append(adbDevices, AndroidDevice{
					DeviceID: parts[0],
					Status:   "connected",
				})
			case "offline":
				adbDevices = append(adbDevices, AndroidDevice{
					DeviceID: parts[0],
					Status:   "offline",
				})
			case "unauthorized":
				adbDevices = append(adbDevices, AndroidDevice{
					DeviceID: parts[0],
					Status:   "unauthorized",
				})
			}
		}
	}
	return adbDevices, nil
}

// ConnectDevice connects to a device via TCP/IP
func ConnectDevice(address string) error {
	adbPath, err := utils.GetADBPath()
	if err != nil {
		return err
	}
	cmd := exec.Command(adbPath, "connect", address)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("adb connect failed: %v, output: %s", err, string(output))
	}
	if strings.Contains(string(output), "unable to connect") || strings.Contains(string(output), "failed to connect") {
		return fmt.Errorf("adb connect failed: %s", string(output))
	}
	return nil
}

// PairDevice pairs with a device using a pairing code
func PairDevice(address, code string) error {
	adbPath, err := utils.GetADBPath()
	if err != nil {
		return err
	}
	cmd := exec.Command(adbPath, "pair", address, code)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("adb pair failed: %v, output: %s", err, string(output))
	}
	if !strings.Contains(string(output), "Successfully paired") {
		return fmt.Errorf("adb pair failed: %s", string(output))
	}
	return nil
}

// GetVideoEncoders returns a list of supported video encoders for a device
func GetVideoEncoders(deviceID string) ([]string, error) {
	adbPath, err := utils.GetADBPath()
	if err != nil {
		return nil, err
	}

	// Query MediaCodecList for video encoders
	cmd := exec.Command(adbPath, "-s", deviceID, "shell", "mediacodec", "-i")
	output, err := cmd.Output()
	if err != nil {
		// Fallback: try dumpsys media.codec
		cmd = exec.Command(adbPath, "-s", deviceID, "shell", "dumpsys", "media.codec", "-v")
		output, err = cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to query encoders: %v", err)
		}
	}

	var encoders []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Look for encoder patterns
		if strings.Contains(line, "Encoder") || strings.Contains(line, "encoder") {
			// Extract encoder name
			parts := strings.Fields(line)
			for _, part := range parts {
				// Skip if it's just "Encoder" keyword
				if strings.EqualFold(part, "Encoder") || strings.EqualFold(part, "encoder") {
					continue
				}
				// Skip if contains punctuation or is too short
				if len(part) < 5 || strings.ContainsAny(part, "(),:") {
					continue
				}
				// Skip common non-encoder keywords
				if strings.Contains(part, "type") || strings.Contains(part, "mime") {
					continue
				}
				// Found an encoder name
				if !contains(encoders, part) {
					encoders = append(encoders, part)
				}
			}
		}
	}

	// Alternative: parse mediacodec -i output more specifically
	if len(encoders) == 0 {
		// Try another approach - look for codec names
		for _, line := range lines {
			// Match patterns like "OMX.qcom.video.encoder.avc" or "c2.android.avc.encoder"
			if strings.Contains(line, ".video.") && (strings.Contains(line, "encoder") || strings.Contains(line, "Encoder")) {
				// Extract the codec name before any space or comma
				fields := strings.Fields(line)
				for _, field := range fields {
					if strings.HasPrefix(field, "OMX.") || strings.HasPrefix(field, "c2.") || strings.HasPrefix(field, "OMX.google.") {
						if !contains(encoders, field) {
							encoders = append(encoders, field)
						}
					}
				}
			}
		}
	}

	// If still empty, return default encoders
	if len(encoders) == 0 {
		encoders = []string{
			"OMX.qcom.video.encoder.avc",
			"c2.android.avc.encoder",
			"OMX.google.h264.encoder",
			"OMX.qcom.video.encoder.hevc",
			"c2.android.hevc.encoder",
		}
	}

	return encoders, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
