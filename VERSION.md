# CiBi Version History

## v0.1.0 (Current Release)
- create `.env` reader
- create `cibi.yaml` reader
- capability to get all file within current dir with prefix file name `*cibi.yaml`
- create documentation how to connect the API to the notion page by add connection
- makes a documentation to prepare the table (database) in notion with basic field:
    - `Name` with type `Title`
    - `Description` with Type `Text`
    - `Integrated with` with Type `Relation` without Limit and related to This database
    - `In Charge` with Type `Person`
    - `Team` with Type `Multi-select`
    - `Repo URL` with Type `Text`
    - `Domain` with Type `URL`
    - `Last edited By`
    - `Last edited Time`
    - `Created By`
- capable to read system environment variable to read basic requirement variable
- capable to read `.env` and read the integrate service
- add basic requirement for `cibi.yaml`
    - key `metadata`
        - key `name`capable to name the services and it will treat as page name in notion
        - key `repoURL` to descibe the repo url
        - key `desc`to descibe the description of the services
        - key `domain` to describe the domain URI of the services (optional and possible to leave it blank)
    - key `readIntegration` for the read the Configuration app
        - key `enable` so this application didn't read the `.env` file to get the integration (default True)
        - key `fileName` File name Config file
        - key `filePath` File path config file
    - key `notionConfig` to enable set database to notion
        - key `apiKey` to define the API Key, use the {{VARIABLE_NAME}} to get this variable from environment file
        - key `databaseID` to define the database ID in notion, use the {{VARIABLE_NAME}} to get this variable from
- add table on every page software with information field
    - Name (APP Deploy Log)
    - deployed_at
    - deployed by
    - domain app
    - relted page
    - team
    - environment
    - start at
    - end at
    - pipeline success status
    - pipeline link
    - github ref
    - triggered by github ref type
    - github run ID (format = github.com/ GITHUB_REPOSITORY)
- https://developers.notion.com/reference/request-limits
- https://docs.github.com/en/github-ae@latest/actions/learn-github-actions/variables#default-environment-variables
- detect if the env contain host and make the integration
- connect to notion through API
- add new record when the service didn't exist
- update the record when the service existed
- create the unit test
- the tutorial to integrate this with github action
- basic setup with 
- example to run on the project
- integrate two of my project to CiBi (one library and one with domain)
- release the docker version so everyone
- release the exec version within release github
- when connected to database remove the secret
- capability to running this locally

## v0.2.0 (Not yet implemented)

- makes a documentation to control the table config in notion
- release to github action to release for secure
- tutorial to runnig  this with bitbucket pipeline
- tutorial to running CiBi with Gitlab Pipeline