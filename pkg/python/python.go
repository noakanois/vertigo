package python

import (
	"fmt"
	"os/exec"


	"log"

)


func PythonGif(shoeid string, folderPath string) {
    err := executePython("gif.py", shoeid, folderPath)
	if err != nil {
		log.Fatalf("error: failed to make GIF: %s", err)
	}
	
}

func executePython(pythonScriptPath string, params ...string) error {
	args := append([]string{"run", "python", pythonScriptPath}, params...)
	cmd := exec.Command("poetry", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(cmd)
		return fmt.Errorf("error: failed to execute command: %s", err)
	}
	fmt.Printf("Output for input %v: %s\n", params, string(output))
	return nil
}