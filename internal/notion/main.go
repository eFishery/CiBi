package notion

import (
	// "os"
	"fmt"
	"log"
	"regexp"
	"context"
	"reflect"
	"strconv"
	"strings"
	"net/url"
	"encoding/json"

	"github.com/jomei/notionapi"
	"github.com/eFishery/CiBi/internal/utils"
)

type NotionConfig struct {
	ApiKey 	   			notionapi.Token	`yaml:"apiKey"`
	DatabaseID 			string 			`yaml:"databaseID"`
	DeploymentLog 		string 			`yaml:"deploymentLog"`
}

type NotionRecord struct {
	Name			string			`title,Name`
	Description		string			`rich_text,Description`
	Domain			string			`rich_text,Domain`
	// Team			[]string		`multi_select,Team`
	IntegratedWith	string			`relation,Integrated With`
	RepoURL			string			`rich_text,Repo URL`
	Environment		string			`select,Environment`
}

type DeploymentLog struct{
	Name			string			`title,Name`
	Domain			string			`rich_text,Domain`
	DeployedAt		notionapi.Date	`date,Deployed at`
	DeployedBy		string			`rich_text,Deployed by`
	Environment		string			`select,Environment`
	PipelineLink	string			`rich_text,Pipeline Link`
	Version			string			`rich_text,Version`
	RepoURL			string			`rich_text,Repo URL`
}

func (nt *NotionConfig) Load(cfg utils.CiBiConfig) *NotionConfig{
	nt.ApiKey = notionapi.Token(cfg.NotionConfig.ApiKey)
	nt.DatabaseID = cfg.NotionConfig.DatabaseID
	nt.DeploymentLog = cfg.NotionConfig.DeploymentLog

	return nt
}


func (nt *NotionConfig) FindDomain(domain string) (*notionapi.Page, error) {
	client := notionapi.NewClient(nt.ApiKey)

	query := &notionapi.DatabaseQueryRequest{ // https://pkg.go.dev/github.com/jomei/notionapi#DatabaseQueryRequest
		Filter: notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: "Domain",
				RichText: &notionapi.TextFilterCondition{
					Contains: domain,
				},
			},
		},
	} //https://github.com/jomei/notionapi/blob/main/database_test.go

	db, err := client.Database.Query(context.Background(), notionapi.DatabaseID(nt.DatabaseID), query) //https://pkg.go.dev/github.com/jomei/notionapi#DatabaseQueryResponse
	if err != nil {
		return nil, err
	}

	if len(db.Results) == 0 {
		return nil, fmt.Errorf("Data Does not Exist for `%s`", domain)
	}

	page, err := client.Page.Get(context.Background(), notionapi.PageID(db.Results[0].ID)) //https://pkg.go.dev/github.com/jomei/notionapi#PageID //https://api.notion.com/v1/pages/:page_id/properties/:property_id
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (nt *NotionConfig) FindRecordDomain(name string, domain string) (*notionapi.Page, error) {
	client := notionapi.NewClient(nt.ApiKey)

	query := &notionapi.DatabaseQueryRequest{ // https://pkg.go.dev/github.com/jomei/notionapi#DatabaseQueryRequest
		Filter: notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: "Name",
				RichText: &notionapi.TextFilterCondition{
					Contains: name,
				},
			},
			notionapi.PropertyFilter{
				Property: "Domain",
				RichText: &notionapi.TextFilterCondition{
					Contains: domain,
				},
			},
		},
	} //https://github.com/jomei/notionapi/blob/main/database_test.go

	db, err := client.Database.Query(context.Background(), notionapi.DatabaseID(nt.DatabaseID), query) //https://pkg.go.dev/github.com/jomei/notionapi#DatabaseQueryResponse
	if err != nil {
		return nil, err
	}

	if len(db.Results) == 0 {
		return nil, fmt.Errorf("Data Does not Exist for `%s` and `%s`", name, domain)
	}

	if len(db.Results) > 1 {
		log.Println("We got", len(db.Results), "results")
		for i, v := range db.Results {
			log.Println("===> Check Page", v.URL)
			pages, err := client.Page.Get(context.Background(), notionapi.PageID(v.ID))
			if err != nil {
				return nil, err
			}
			FieldToCompare := []string{"Name", "Domain"}
			dataResult := []string{}
			for _, vv := range FieldToCompare {
				getValue, err := nt.PropertiesValue(pages, vv)
				if err != nil {
					return nil, err
				}
				dataResult = append(dataResult, getValue[0])
			}
			if dataResult[0] == name && dataResult[1] == domain {
				log.Printf("This page have a same value with Name `%s` and Domain `%s`", dataResult[0], dataResult[1])
				return pages, nil
			}

			log.Printf("This page have a different value with Name `%s` and Domain `%s`", dataResult[0], dataResult[1])

			if (i + 1) == len(db.Results) {
				return nil, fmt.Errorf("Data Does not Exist for `%s` and `%s`", name, domain)
			}
		}
	}

	page, err := client.Page.Get(context.Background(), notionapi.PageID(db.Results[0].ID)) //https://pkg.go.dev/github.com/jomei/notionapi#PageID //https://api.notion.com/v1/pages/:page_id/properties/:property_id
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (nt *NotionConfig) GetPageTitle(PageID string) (string, error) {
	client := notionapi.NewClient(nt.ApiKey)
	page, err := client.Page.Get(context.Background(), notionapi.PageID(PageID))
	if err != nil {
		return "", err
	}

	res2B, _ := json.Marshal(page.Properties["Name"])
	var dat map[string]interface{}
	err = json.Unmarshal(res2B, &dat)
	if err != nil {
		return "", err
	}
	return dat["title"].([]interface{})[0].(map[string]interface{})["plain_text"].(string), nil
}

// func (nt *NotionConfig) NewDatabase(tableName string) error {
// 	client := notionapi.NewClient(nt.ApiKey)

// 	dataSet := &notionapi.DatabaseCreateRequest{
// 		Parent: notionapi.Parent{
// 			Type:   notionapi.ParentTypePageID,
// 			PageID: notionapi.PageID("cf18546d662a4da38ccfcea47bae74a8"),
// 			// DatabaseID: notionapi.DatabaseID(utils.Settings.NotionMappingPage),
// 		},
// 		Title: []notionapi.RichText{
// 			{
// 				Type: notionapi.ObjectTypeText,
// 				Text: notionapi.Text{Content: tableName},
// 			},
// 		},
// 		Properties: notionapi.PropertyConfigs{
// 			"Name": notionapi.TitlePropertyConfig{
// 				Type: notionapi.PropertyConfigTypeTitle,
// 			},
// 		},
// 	}

// 	_, err := client.Database.Create(context.Background(), dataSet)

// 	if err != nil {
// 		return fmt.Errorf("Create() error = %v", err)
// 	}

// 	return nil
// }

// func (nt *NotionConfig) UpdateRecord(Page *notionapi.Page, mapping utils.MappingKey, newValue string) error {
// 	utils.Settings = utils.LoadSetting()
// 	client := notionapi.NewClient(utils.Settings.NotionKey)

// 	dataSet := &notionapi.PageUpdateRequest{
// 		Properties: notionapi.Properties{},
// 	}

// 	switch mapping.FieldType {
// 	case "title":
// 		dataSet.Properties[mapping.FieldName] = notionapi.TitleProperty{
// 			Title: []notionapi.RichText{
// 				{Text: notionapi.Text{Content: newValue}},
// 			},
// 		}
// 	case "rich_text":
// 		dataSet.Properties[mapping.FieldName] = notionapi.RichTextProperty{
// 			RichText: []notionapi.RichText{
// 				{Text: notionapi.Text{Content: newValue}},
// 			},
// 		}
// 	case "checkbox":
// 		dataSet.Properties[mapping.FieldName] = notionapi.CheckboxProperty{
// 			Checkbox: true,
// 		}
// 	case "multi_select":
// 		dataSet.Properties[mapping.FieldName] = notionapi.MultiSelectProperty{
// 			MultiSelect: []notionapi.Option{
// 				{Name: newValue},
// 			},
// 		}
// 	}

// 	_, err := client.Page.Update(context.Background(), notionapi.PageID(Page.ID), dataSet)

// 	if err != nil {
// 		return fmt.Errorf("Create() error = %v", err)
// 	}

// 	return nil

// }

func (nt *NotionConfig) UpdateRecordRelation(Page *notionapi.Page, newValue []notionapi.PageID) error {
	client := notionapi.NewClient(nt.ApiKey)

	dataSet := &notionapi.PageUpdateRequest{
		Properties: notionapi.Properties{},
	}

	nV := []notionapi.Relation{}

	for _, val := range newValue {
		nV = append(nV, notionapi.Relation{ID: val})
	}

	dataSet.Properties["Integrated with"] = notionapi.RelationProperty{
		Relation: nV,
	}
	_, err := client.Page.Update(context.Background(), notionapi.PageID(Page.ID), dataSet)

	if err != nil {
		return fmt.Errorf("Create() error = %v", err)
	}

	return nil

}

func (nt *NotionConfig) NewRecord(dataRecord *NotionRecord) (*notionapi.Page, error) {
	client := notionapi.NewClient(nt.ApiKey)

	dataSet := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(nt.DatabaseID),
		},
		Properties: notionapi.Properties{},
	}

	structNC := reflect.ValueOf(dataRecord).Elem()
	for i := 0; i < structNC.NumField(); i++ {
		fieldNC := structNC.Type().Field(i).Name
		fieldValue := structNC.Field(i).Interface()
		fieldTag, _ := reflect.TypeOf(dataRecord).Elem().FieldByName(fieldNC)
		fieldType := strings.Split(string(fieldTag.Tag), ",")[0]
		fieldName := strings.Split(string(fieldTag.Tag), ",")[1]
		switch fieldType {
		case "title":
			dataSet.Properties[fieldName] = notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{Text: &notionapi.Text{
						Content: fmt.Sprintf("%v", fieldValue),},
					},
				},
			}
		case "rich_text":
			dataSet.Properties[fieldName] = notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{Text: &notionapi.Text{
						Content: fmt.Sprintf("%v", fieldValue),},
					},
				},
			}
		case "multi_select":
			dataSet.Properties[fieldName] = notionapi.MultiSelectProperty{
				MultiSelect: []notionapi.Option{
					{Name: fmt.Sprintf("%v", fieldValue),},
				},
			}
		case "select":
			dataSet.Properties[fieldName] = notionapi.SelectProperty{
				Select: notionapi.Option{
					Name: fmt.Sprintf("%v", fieldValue),
				},
			}
		case "date":
			var dataField *notionapi.Date
			dataField = fieldValue.(*notionapi.Date)
			dataSet.Properties[fieldName] = notionapi.DateProperty{
				Date: &notionapi.DateObject{
					Start: dataField,
				},
			}
		// can't continue this, and must inputted manually
		// because notion API won't get the "Guest" user
		// case "people":
		// 	var people []notionapi.User
		// 	for _, email := range dataRecord.InCharge {
		// 		people = append(people, notionapi.User{Person: &notionapi.Person{Email: email}})
		// 	}
		// 	dataSet.Properties[fieldName] = notionapi.PeopleProperty{
		// 		People: people,
		// 	}
		}
	}

	pageReturn, err := client.Page.Create(context.Background(), dataSet)

	if err != nil {
		return nil, fmt.Errorf("Create() error = %v", err)
	}

	return pageReturn, nil
}

func (nt *NotionConfig) NewLogRecord(dataRecord *DeploymentLog) (*notionapi.Page, error) {
	client := notionapi.NewClient(nt.ApiKey)

	dataSet := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(nt.DeploymentLog),
		},
		Properties: notionapi.Properties{},
	}

	structNC := reflect.ValueOf(dataRecord).Elem()
	for i := 0; i < structNC.NumField(); i++ {
		fieldNC := structNC.Type().Field(i).Name
		fieldValue := structNC.Field(i).Interface()
		fieldTag, _ := reflect.TypeOf(dataRecord).Elem().FieldByName(fieldNC)
		fieldType := strings.Split(string(fieldTag.Tag), ",")[0]
		fieldName := strings.Split(string(fieldTag.Tag), ",")[1]
		switch fieldType {
		case "title":
			dataSet.Properties[fieldName] = notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{Text: &notionapi.Text{
						Content: fmt.Sprintf("%v", fieldValue),},
					},
				},
			}
		case "rich_text":
			dataSet.Properties[fieldName] = notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{Text: &notionapi.Text{
						Content: fmt.Sprintf("%v", fieldValue),},
					},
				},
			}
		case "multi_select":
			dataSet.Properties[fieldName] = notionapi.MultiSelectProperty{
				MultiSelect: []notionapi.Option{
					{Name: fmt.Sprintf("%v", fieldValue),},
				},
			}
		case "select":
			dataSet.Properties[fieldName] = notionapi.SelectProperty{
				Select: notionapi.Option{
					Name: fmt.Sprintf("%v", fieldValue),
				},
			}
		case "date":
			var dataField notionapi.Date
			dataField = fieldValue.(notionapi.Date)

			dataSet.Properties[fieldName] = notionapi.DateProperty{
				Date: &notionapi.DateObject{
					Start: &dataField,
				},
			}
		// can't continue this, and must inputted manually
		// because notion API won't get the "Guest" user
		// case "people":
		// 	var people []notionapi.User
		// 	for _, email := range dataRecord.InCharge {
		// 		people = append(people, notionapi.User{Person: &notionapi.Person{Email: email}})
		// 	}
		// 	dataSet.Properties[fieldName] = notionapi.PeopleProperty{
		// 		People: people,
		// 	}
		}
	}

	pageReturn, err := client.Page.Create(context.Background(), dataSet)

	if err != nil {
		return nil, fmt.Errorf("Create() error = %v", err)
	}

	return pageReturn, nil
}

// can't do get list users because this isn't available for "Guest" member
// while the objective of this project is to make service mapping free available for everyone to user
func (nt *NotionConfig) GetListUsers() ([]notionapi.User, error){
	client := notionapi.NewClient(nt.ApiKey)

	got, err := client.User.List(context.Background(), &notionapi.Pagination{PageSize: 100})
	if (err != nil) {
		return nil, err
	}

	return got.Results, nil
}

func (nt *NotionConfig) PropertiesValue(Page *notionapi.Page, PropertiesName string) ([]string, error) {
	res2B, err := json.Marshal(Page.Properties[PropertiesName])
	if err != nil {
		return []string{"exception"}, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(res2B, &data)
	if err != nil {
		return []string{"exception"}, err
	}

	keys := make([]string, len(data))
	i := 0
	for k := range data {
		keys[i] = k
		i++
	}

	tipe := data["type"].(string)
	switch reflect.TypeOf(data[tipe]).Name() {
	case "string": // url, create_at
		return []string{data[tipe].(string)}, nil
	case "bool": // checkbox
		return []string{strconv.FormatBool(data[tipe].(bool))}, nil
	case "float64": // number
		return []string{fmt.Sprintf("%.0f", data[tipe])}, nil
	}
	switch tipe {
	case "created_by", "date", "last_edited_by":
		if tipe == "date" {
			dataDate := data[tipe].(map[string]interface{})
			return []string{fmt.Sprintf("%#v", dataDate["start"]), fmt.Sprintf("%#v", dataDate["end"])}, nil
		}
		return []string{data[tipe].(map[string]interface{})["name"].(string)}, nil
	default:
		getValue := data[tipe].([]interface{})
		output := make([]string, len(getValue))
		for index := range getValue {
			if tipe == "title" || tipe == "rich_text" {
				output[index] = getValue[index].(map[string]interface{})["text"].(map[string]interface{})["content"].(string)
			} else if tipe == "people" || tipe == "multi_select" {
				output[index] = getValue[index].(map[string]interface{})["name"].(string)
			} else if tipe == "relation" {
				output[index], err = nt.GetPageTitle(getValue[index].(map[string]interface{})["id"].(string))
				if err != nil {
					return nil, err
				}
			} else {
				fmt.Println("The type ", tipe, " is not registered, with content", data[tipe])
			}
		}
		return output, nil
	}

}

func ExtractDomain(data map[string]string) []string {
	var values []string

	keyword := []string{"http://", "https://", "service.local"}
	for _, v := range data {
		for _, kw := range keyword {
			if strings.Contains(v, kw) {
				values = append(values, v)
				break
			}
		}
	}

	return values
}

func ReplaceRegexp(name string) string {

	return ""
}

func (nt *NotionConfig) SimpleNewRecord(domain string) (*notionapi.Page, error) {
	client := notionapi.NewClient(nt.ApiKey)

	url, err := url.Parse(domain)
	if err != nil {
		log.Fatal(err)
	}

	dataSet := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(nt.DatabaseID),
		},
		Properties: notionapi.Properties{
			"Name": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{Text: &notionapi.Text{Content: url.Hostname()}},
				},
			},

			"Domain": notionapi.RichTextProperty{
				RichText: []notionapi.RichText{
					{Text: &notionapi.Text{Content: regexp.MustCompile(`//.*@`).ReplaceAllStringFunc(domain, ReplaceRegexp)}},
				},
			},
		},
	}
	createNewPage, err := client.Page.Create(context.Background(), dataSet)

	if err != nil {
		return nil, fmt.Errorf("Create() error = %v", err)
	}

	return createNewPage, nil
}

func (nt *NotionConfig) CheckDomainRelation(Page *notionapi.Page, data map[string]string) error {
	DomainList := ExtractDomain(data)

	listPageID := []notionapi.PageID{}

	for _, val := range DomainList {
		// continue the loop if the string is email
		_, isEmail := utils.ValidMailAddress(val)
		if isEmail {
			continue
		}
		result, err := nt.FindDomain(regexp.MustCompile(`//.*@`).ReplaceAllStringFunc(val, ReplaceRegexp))
		if err != nil {
			newPage, err := nt.SimpleNewRecord(val)
			if err != nil {
				return fmt.Errorf("Fail create new record for update relation")
			}
			listPageID = append(listPageID, notionapi.PageID(newPage.ID))
			continue
		}
		_, err = nt.PropertiesValue(result, "Name")
		if err != nil {
			return err
		}
		listPageID = append(listPageID, notionapi.PageID(result.ID))
	}

	err := nt.UpdateRecordRelation(Page, listPageID)
	if err != nil {
		return fmt.Errorf("ERROR UPDATE Record Relation = %v", err)
	}

	return nil
}
