package notionlog 

import (
	"os"
	"fmt"
	"time"

	"github.com/eFishery/CiBi/internal/utils"
	"github.com/eFishery/CiBi/internal/notion"

	"github.com/jomei/notionapi"
	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func RunDeployedLog(cmd *cobra.Command, args []string) {

	var cibiConfig string 

	if len(viper.GetString("file")) > 0 {
		cibiConfig = viper.GetString("file")
	}else{
		cibiConfig = "cibi.yaml"
	}

	var cfg utils.CiBiConfig
	cibiFile := fmt.Sprintf("%s/%s", utils.GetCurrDir(), cibiConfig)
	_, err := cfg.Load(cibiFile)
	if err != nil{
		fmt.Println(" └─[-] File",cibiFile,"is not found, Please create the CiBi config file")
		os.Exit(1)
	}
	var nt notion.NotionConfig
	nt.Load(cfg)

	now := notionapi.Date(time.Now())

	dataRecord := &notion.DeploymentLog{
		Name: cfg.Metadata.Name,
		DeployedAt: now,
		DeployedBy: "yahya.kimochi@gmail.com",
		Domain: cfg.Metadata.Domain,
		RepoURL: cfg.Metadata.RepoURL,
		Environment: cfg.Metadata.Environment,
		PipelineLink: "https://github.com/k1m0ch1/jemawa-menti-choices-spammer/actions/runs/6578783742",
		Version: "GITHUB_REF",
	}

	page, err := nt.NewLogRecord(dataRecord)
	if err != nil {
		fmt.Printf(" └─[-] Exit Application %v\n", err)
		os.Exit(1)
	}

	fmt.Println(page.ID)
}