TAIGO Examples & Demo Cases
-----

![alt text](../assets/banner_627x300.png "TAIGO Banner")

# Introduction
Each section in this document is a use case snippet. These snippets are useless as-is, because they are taken our from a grater context.

Make sure you always import the required packages and you always put these snippets into the right context.

There is no proper error handling implemented in these snippets, each error panics the application. Make sure you always implement proper error handling.

# Cases

## Create Client & Authenticate
```go
import (
    "fmt"
    "net/http"

    taiga "github.com/theriverman/taigo"
)

func main() {
	// Create client
	client := taiga.Client{
		BaseURL:    "https://api.taiga.io",
		HTTPClient: &http.Client{},
	}
	// Initialise client (authenticates to Taiga)
	err := client.Initialise()
	if err != nil {
		panic(err)
	}

	// Authenticate (get/set Token)
	client.AuthByCredentials(&taiga.Credentials{
		Type:     "normal",
		Username: "admin",
		Password: "123123",
	})
}

```

## Get Self (User) /users/me
```go
me, err := client.User.Me()
if err != nil {
	panic(err)
}
fmt.Println("Me: (ID, Username, FullName)", me.ID, me.Username, me.FullName)
```

## Get a Project
```go
// Get Project
slug := "my-sassy-project-1"
fmt.Printf("\nGetting Project (slug=%s)..\n", slug)
project, err := client.Project.GetBySlug(slug)
if err != nil {
	panic(err)
}
fmt.Printf("\nProject name: %s \n", project.Name)
```

## Get Project Severities
In this scenario observe that we're accessing the `ProjectDETAIL` meta field.
```go
fmt.Printf("\nProject Severities:\n")
for _, severity := range project.ProjectDETAIL.Severities {
	fmt.Printf("  * ID=%d Name=%s\n", severity.ID, severity.Name)
}
```

## Get (Project) Epic Custom Attributes
A total of 3 epic custom attributes are returned. See it limited by the for-loop.
```go
fmt.Println("\nGetting all (Project) Epics Custom Attributes and printing the first 3 to the console:")
fmt.Printf("Project Epic Custom Attributes:\n")
for _, epicCA := range project.ProjectDETAIL.EpicCustomAttributes {
	fmt.Printf("  * ID=%d Name=%s ProjectID=%d\n", epicCA.ID, epicCA.Name, epicCA.ProjectID)
}
```

## Get Epics
A total of 3 epics are returned. See it limited by the for-loop.
```go
fmt.Println("\nGetting all Epics and printing the first 3 to the console:")
epics, err := client.Epic.List(nil)
if err != nil {
	panic(err)
}
for i := 0; i < 3; i++ {
	epic := epics[i]
	fmt.Printf("  * epics[%d] :: ID=%d | Subject=%s\n", i, epic.ID, epic.Subject)
	// Access ModifiedDate via meta because that's not available through the generic `Epic`
	meta := *epic.EpicDetailLIST
	fmt.Printf("  * meta :: ModifiedDate = %s\n\n", meta[i].ModifiedDate.Format("2006-01-02 15:04:05"))
}
```

## Get Epic by ID
```go
epicID := 123456
fmt.Println("Getting an Epic by ID:", epicID)
epic, err := client.Epic.Get(epicID)
if err != nil {
    panic(err)
}
fmt.Println("  * epic.EpicDetailGET.ID", epic.EpicDetailGET.ID)
fmt.Printf("  * epic.Subject = %s\n\n", epic.Subject)
```

## Get User Story Custom Attributes
```go
storyID := 123456
cavs := taiga.UserStoryCustomAttribValues{}
_, err = p.client.Request.Get(p.client.MakeURL("userstories", "custom-attributes-values", strconv.Itoa(storyID)), &cavs)
if err != nil {
    panic(err)
}
for attrDefId, attrValue := range cavs.AttributesValues {
    fmt.Printf("  * DefinitionID=%s Value=%v\n", attrDefId, attrValue)
}
```

## Update User Story Custom Attributes
```go
storyID := 123456
// All fields have to be sent: https://github.com/taigaio/taiga-back/issues/1029
updateAttributes := taiga.TgObjCAVDBase{}
updateAttributes.Version = cavs.Version // copy old/original version number (see previous example)
updateAttributes.AttributesValues = cavs.AttributesValues // copy old/original attribute values (see previous example)
// Now override attribute values that should be changed
updateAttributes.AttributesValues["1"] = "value1"
updateAttributes.AttributesValues["2"] = 987654321
_, err = p.client.Request.Patch(p.client.MakeURL("userstories", "custom-attributes-values", strconv.Itoa(storyID)), updateAttributes, nil)
if err != nil {
    panic(err)
}
```