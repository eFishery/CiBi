package servicemap

import(
	"os"
	"fmt"

	"github.com/eFishery/CiBi/internal/utils"
	"github.com/eFishery/CiBi/internal/notion"

	godot "github.com/joho/godotenv"
	cobra "github.com/spf13/cobra"
)

func RunServiceMap(cmd *cobra.Command, args []string) {

	fmt.Println("[+] Running CiBi for Service Mapping")

	var cibiConfig string

	if len(args) > 0 {
		cibiConfig = args[0]
		fmt.Println("[+] Trying to load", cibiConfig)
	}else{
		fmt.Println("[+] CiBi file didn't set up, will used cibi.yaml as default")
		cibiConfig = "cibi.yaml"
	}

	var cfg utils.CiBiConfig
	cibiFile := fmt.Sprintf("%s/%s", utils.GetCurrDir(), cibiConfig)
	_, err := cfg.Load(cibiFile)
	if err != nil{
		fmt.Printf(" └─[-] %v\n", err)
		fmt.Println(" └─[-]",cibiFile,"Please create the CiBi config file")
		os.Exit(1)
	}
	fmt.Println(" ├─[+]", cibiFile, "is Loaded")

	fmt.Printf(" ├─[+] Register service mapping for `%s`\n", cfg.Metadata.Name)

	// check if the service is registered into notion

	var nt notion.NotionConfig
	nt.Load(cfg)
	pg, err := nt.FindRecordDomain(cfg.Metadata.Name, cfg.Metadata.Domain)
	if err != nil {
		fmt.Printf(" |  ├─[-] Page `%s` Does not exist\n", cfg.Metadata.Name)
		fmt.Printf(" |  |  ├─[-] Error message : %v\n", err)
		fmt.Printf(" |  |  ├─[+] Create new Page called `%s`\n", cfg.Metadata.Name)

		// if not exist registered create one
		dataRecord := &notion.NotionRecord{
			Name: cfg.Metadata.Name,
			Description: cfg.Metadata.Description,
			Domain: cfg.Metadata.Domain,
			RepoURL: cfg.Metadata.RepoURL,
			Environment: cfg.Metadata.Environment,
		}

		page, err := nt.NewRecord(dataRecord)
		if err != nil {
			fmt.Printf(" |  |  |  ├─[-] Error message : %v\n", err)
			fmt.Printf(" |  |  |  └─[-] Exit Application : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf(" |  |  └─[+] Page Created with URL %s\n", page.URL)

	}else{
		fmt.Printf(" |  ├─[+] Page `%s` already exist at %s, proceed to update value\n", cfg.Metadata.Name, pg.URL)
	}

	// fill up the integration with field
	if cfg.ReadIntegration.Enable{
		fmt.Printf(" |  └─[+] Read Integration is Enabled, will read file %s%s\n", cfg.ReadIntegration.FilePath, cfg.ReadIntegration.FileName)
		file := fmt.Sprintf("%s%s", cfg.ReadIntegration.FilePath, cfg.ReadIntegration.FileName)

		var rIEnv map[string]string
		rIEnv, err = godot.Read(file)
		if err != nil {
			fmt.Printf(" |  └─[-] %v\n", err)
			fmt.Println(" |  └─[-]",file,"can't be opened")
			os.Exit(1)
		}

		page, err := nt.FindRecordDomain(cfg.Metadata.Name, cfg.Metadata.Domain)
		if err != nil {
			fmt.Printf(" |  |  └─[-] Error when Find Page %s&%s\n", cfg.Metadata.Name, cfg.Metadata.Domain, err)
			fmt.Printf(" |  |  └─[-] with error Message = %v")
		}
		err = nt.CheckDomainRelation(page, rIEnv)
		if err != nil {
			fmt.Printf(" |  |  └─[-] %v\n", err)
		}
		fmt.Printf(" |  └─[+] End of Register `%s`\n", cfg.Metadata.Name)
	}

	fmt.Println(" └─[+] End of Program")
}