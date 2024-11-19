package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/model"
)




func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", model.NewCustomError(500, fmt.Sprintf("command failed: %v\nOutput: %s", err, out.String()))
	}
	return out.String(), nil
}
