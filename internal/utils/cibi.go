package utils

type Metadata struct {
	Name		 	string    		`yaml:"name"`
	Description    	string    		`yaml:"description"`
	Domain       	string    		`yaml:"domain"`
	RepoURL			string			`yaml:"repoURL"`
	Team			[]string		`yaml:"team"`
	Environment		string			`yaml:"environment"`
}

type ReadIntegration struct {
	Enable   		bool   			`yaml:"enable"`
	FileName 		string 			`yaml:"fileName"`
	FilePath 		string   		`yaml:"filePath"`
}

type NotionConfig struct {
	ApiKey 	   		string			`yaml:"apiKey"`
	DatabaseID 		string 			`yaml:"databaseID"`
	DeploymentLog	string			`yaml:"deploymentLog"`
}

type CiBiConfig struct{
	Metadata		Metadata		`yaml:"metadata"`
	ReadIntegration	ReadIntegration	`yaml:"readIntegration"`
	NotionConfig	NotionConfig	`yaml:"notionConfig"`
	FileName		string
}