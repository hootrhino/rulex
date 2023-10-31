package ossupport

import (
	"fmt"
	"os"
	"os/exec"
)

/*
*
* 以后可能用于非 systemctl
*
 */
func SendKillSignal(pid int) error {
	cmd := exec.Command("kill", "-9", fmt.Sprintf("%d", pid))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
func GetProcessStatus(pid int) (string, error) {
	cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", pid))
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	output := string(cmdOutput)
	return output, nil
}
