package test

// import (
// 	"gin-template-app/pkg/log"
// 	"gin-template-app/pkg/utils"
// 	"os"
// 	"os/exec"
// 	"testing"
// )

// // Data data struct
// type Data = map[string]interface{}

// // IDS ids struct
// type IDS = []uint

// // InitTests Init tests module
// func InitTests(t *testing.T) {
// 	t.Helper()
// 	// get config file in default path: config/config.yaml
// 	thisPath, err := os.Getwd()
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	source := utils.GetSubPath(thisPath, 2) + "/config/"
// 	log.Infoln(source)
// 	// copy file to the default path
// 	if err = exec.Command("mkdir", "config").Run(); err != nil {
// 		log.Panicln(err)
// 	}
// 	if err = exec.Command("cp", "-r", source, ".").Run(); err != nil {
// 		log.Panicln(err)
// 	}

// 	// setup config
// 	if _, err = setting.SetupConfig(); err != nil {
// 		log.Panicln(err)
// 	}

// 	// delete config
// 	if err = exec.Command("rm", "-rf", "config/").Run(); err != nil {
// 		log.Panicln(err)
// 	}
// }
