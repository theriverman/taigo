package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	taiga "github.com/theriverman/taigo"
)

var sandboxEpicID2 int // do not change manually

/*
  * Before running this main.go file, please read the whole code below!
  * Pay attention to the inline/block comments and plan your Taiga sandbox project ahead!

  * To successfully run this example, you must have:
	* A sandbox Taiga project created -> client.SetDefaultProjectBySlug()
	* Have 3 epics
	* Have 3 milestones (sprints)
	* Have 3 User Stories
	* Have 3 Tasks
*/

// Set via build flag (see build scripts)
var taigaHost = "https://api.taiga.io"
var taigaAuthType = "normal"
var taigaUsername string
var taigaPassword string

// Project Details
var sandboxProjectSlug string
var sandboxEpicID string
var sandboxFileUploadPath string // absolute path your file to be uploaded

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetPrefix("TAIGO DEV" + " >> ")

	/**
	 * Uncomment and initialize the following variables if you're not providing
	 * the variable values through build flags (see the build-and-run scripts)
	 */

	// sandboxProjectSlug = ""
	// sandboxEpicID = ""
	// sandboxFileUploadPath = ""
	/**
	 * EXAMPLES:
	 * fileAttachmentUploadPath = "/home/user/Documents/Bad-puns-make-me-sic.1.jpg"	// LINUX
	 * fileAttachmentUploadPath = "C:\\Images\\Bad-puns-make-me-sic.1.jpg"				// WINDOWS
	 */

	// Convert Epic ID string to int
	var e error
	sandboxEpicID2, e = strconv.Atoi(sandboxEpicID)
	if e != nil {
		log.Fatalln(e)
		return
	}
}

func main() {
	// Create client
	client := taiga.Client{
		BaseURL:    taigaHost,
		HTTPClient: &http.Client{},
	}
	// Initialise client (authenticates to Taiga)
	err := client.Initialise()
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Authenticate (get/set Token)
	err = client.AuthByCredentials(&taiga.Credentials{
		Type:     taigaAuthType,
		Username: taigaUsername,
		Password: taigaPassword,
	})
	if err != nil {
		log.Println("Error!", err)
		return
	}

	// Set default project (optional. recommended for convenience)
	client.SetDefaultProjectBySlug(sandboxProjectSlug)

	// Get /users/me
	me, err := client.User.Me()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Me: (ID, Username, FullName, DateJoined)", me.ID, me.Username, me.FullName, me.DateJoined.String())

	// Get Project
	fmt.Printf("\nGetting Project (slug=%s)..\n", sandboxProjectSlug)
	project, err := client.Project.GetBySlug(sandboxProjectSlug)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("\nProject name: %s | ID:%d \n", project.Name, project.ID)

	// Get Project Agile Points
	fmt.Println("Registered Agile Points for Project:", project.Name)
	for _, agilepoint := range project.ProjectDETAIL.Points {
		if !agilepoint.IsValueNil() {
			fmt.Println("agilepoint.Value ->", *agilepoint.Value)
		}
	}

	// Get Project Severities
	fmt.Printf("\nProject Severities:\n")
	for _, severity := range project.ProjectDETAIL.Severities {
		fmt.Printf("  * ID=%d Name=%s\n", severity.ID, severity.Name)
	}

	// Get Project Epic Custom Attributes
	log.Println("\nGetting all (Project) Epics Custom Attributes and printing the first 3 to the console:")
	fmt.Printf("Project Epic Custom Attributes:\n")
	for _, epicCA := range project.ProjectDETAIL.EpicCustomAttributes {
		fmt.Printf("  * ID=%d Name=%s ProjectID=%d\n", epicCA.ID, epicCA.Name, epicCA.ProjectID)
	}

	// Get Epics
	// (total of 3; limited by the for-loop)
	log.Println("\nGetting all Epics and printing the first 3 to the console:")
	epics, err := client.Epic.List(nil)
	if err != nil {
		log.Println(err)
		return
	}
	for i := 0; i < 3; i++ {
		epic := epics[i]
		fmt.Printf("  * epics[%d] :: ID=%d | Subject=%s\n", i, epic.ID, epic.Subject)
		// Accessing ModifiedDate via meta because that's not available through the generic `Epic`
		meta := *epic.EpicDetailLIST
		fmt.Printf("  * meta :: ModifiedDate = %s\n\n", meta[i].ModifiedDate.Format("2006-01-02 15:04:05"))
	}

	// Get Epic by ID
	log.Println("Getting an Epic by ID:", sandboxEpicID2)
	epic, err := client.Epic.Get(sandboxEpicID2)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("  * epic.EpicDetailGET.ID || epic.ID", epic.EpicDetailGET.ID, epic.ID)
	fmt.Printf("  * epic.Subject = %s\n\n", epic.Subject)

	// Get milestones (for default project if set)
	// (total of 3; limited by the for-loop)
	log.Println("Getting all Milestones(Sprints) and printing the first 3 to the console:")
	milestones, mti, err := client.Milestone.List(nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("  *  Taiga-Info-Total-Opened-Milestones =", mti.TaigaInfoTotalOpenedMilestones)
	log.Println("  *  Taiga-Info-Total-Closed-Milestones =", mti.TaigaInfoTotalClosedMilestones)
	if len(milestones) >= 3 {
		for i := 0; i < 3; i++ {
			milestone := milestones[i]
			log.Println("  * ", milestone.ID, milestone.Name)
		}
	}

	// Get UserStories ( for project.ID -- see UserStoryQueryParams )
	// (total of 3; limited by the for-loop)
	log.Println("\nGetting all UserStories and printing the first 3 to the console:")
	userstories, err := client.UserStory.List(
		&taiga.UserStoryQueryParams{
			// IncludeAttachments: true,  // set to true to include attachments
			Project: project.ID,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	for i := 0; i < 3; i++ {
		us := userstories[i]
		meta := (*us.UserStoryDetailLIST)[i] // Dereference to the MetaList, then access by slice index
		log.Println("  * ", us.ID, us.Subject)
		log.Println("  *  Meta :: ID", meta.ID)
		for _, attachment := range meta.Attachments {
			fmt.Println(attachment.ThumbnailCardURL)
		}
	}

	/*
		RAW REQUEST
		List ALL Epic custom attributes by composing a RAW request
	*/
	log.Println("Getting all Epic custom attributes and printing the first 3 to the console:")

	/*
		We declare the EpicCustomAttributeDetail struct here so the returned JSON payload
		can be serialized into this struct
	*/

	// EpicCustomAttributeDetail -> https://taigaio.github.io/taiga-doc/dist/api.html#object-epic-custom-attribute-detail
	// Converted via https://mholt.github.io/json-to-go/
	type EpicCustomAttributeDetail struct {
		CreatedDate  time.Time   `json:"created_date"`
		Description  string      `json:"description"`
		Extra        interface{} `json:"extra"`
		ID           int         `json:"id"`
		ModifiedDate time.Time   `json:"modified_date"`
		Name         string      `json:"name"`
		Order        int         `json:"order"`
		Project      int         `json:"project"`
		Type         string      `json:"type"`
	}
	epicCustomAttributes := []EpicCustomAttributeDetail{} // this will hold the returned data. passed as a pointer below

	// Final URL will be be https://*******/api/v1/epic-custom-attributes
	resp, err := client.Request.Get(client.MakeURL("epic-custom-attributes"), &epicCustomAttributes)
	if err != nil {
		log.Println(err)
		log.Println(resp)
		return
	}
	for i := 0; i < 3; i++ {
		ca := epicCustomAttributes[i]
		log.Println("  * ", ca.ID, ca.Name)
	}

	/*
		CREATE OPERATIONS
	*/

	// Create a new Epic. Subject = My Epic @ `time.Now()`
	log.Println("\nCreating a new Epic and printing its ID and Subject to console:")
	newOutEpic := taiga.Epic{
		Subject: fmt.Sprintf("My Epic @ %s", time.Now().String()),
		Project: project.ID,
	}
	newEpic, err := client.Epic.Create(&newOutEpic)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("  * ", newEpic.ID, newEpic.Subject)

	// Create a new Milestone (Sprint). Subject = My Sprint @ `time.Now()`
	log.Println("\nCreating a Sprint and printing its ID and Name to console:")
	newSprint := taiga.Milestone{
		Name:            fmt.Sprintf("My Sprint @ %s", time.Now().String()),
		Project:         project.ID,
		EstimatedStart:  "2020-05-17",
		EstimatedFinish: "2020-05-27",
	}
	sprint, err := client.Milestone.Create(&newSprint)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("  * ", sprint.ID, sprint.Name)

	// Create a new UserStory. Subject = My US @ `time.Now()`
	/*
		This is an alternative approach. Here the newly created UserStory details
		are fed back into the original UserStory struct:
		  * See how `us := &taiga.UserStory{}` is now a pointer variable
		  * See how `client.UserStory.Create(us)` now takes a US pointer
		  * See how `us, err = client.UserStory.Create(*us)` feeds response into the original US
	*/
	log.Println("\nCreating a new UserStory and printing its ID and Subject to console:")
	us := &taiga.UserStory{
		Subject:   fmt.Sprintf("My US @ %s", time.Now().String()),
		Project:   project.ID,
		Milestone: sprint.ID,
	}
	us, err = client.UserStory.Create(us)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("  * ", us.ID, us.Subject)

	// Relate the UserStory to an Epic
	fmt.Printf("\nRelating UserStory (ID=%d) to Epic (ID=%d) console:\n", us.ID, newEpic.ID)
	_, err = us.RelateToEpic(&client, epic.ID)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("\n  * UserStory (ID=%d) has been related to Epic (ID=%d)\n", us.ID, newEpic.ID)

	// Create Task for UserStory
	fmt.Printf("\nAdding a Task(Subject=My Task @ `time.Now()) to UserStory (ID=%d):\n", us.ID)
	task := &taiga.Task{
		Subject:   fmt.Sprintf("My Task @ %s", time.Now().String()),
		Project:   project.ID,
		UserStory: us.ID,
	}
	task, err = us.CreateRelatedTask(&client, *task)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("\n  * Task (ID=%d) has been created\n", task.ID)

	// Upload a file to our Task
	fmt.Printf("\nAdding an Attachment to Task(Subject=My Task @ `time.Now()):\n")
	attachment := &taiga.Attachment{
		Name:        "My Fancy Avatar",
		Description: "This is a test file uploaded via TAIGO",
	}
	attachment.SetFilePath(sandboxFileUploadPath)
	attachment, err = client.Task.CreateAttachment(attachment, task)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("\n  * Task (ID=%d) Attachment (ID=%d) has been created\n", task.ID, attachment.ID)

	// Get and check out just uploaded Task Attachment
	fmt.Printf("\nGetting Task Attachment (Subject=My Task @ `time.Now()):\n")
	attachmentAgain, err := client.Task.GetAttachment(attachment.ID)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("\n  * Attachment (ID=%d) has size: %d\n", attachmentAgain.ID, attachmentAgain.Size)

	// List Issues
	log.Println("Getting all issues [for default project]")
	issueList, err := client.Issue.List(&taiga.IssueQueryParams{Project: client.GetDefaultProjectID()})
	if err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < 3; i++ {
		issue := issueList[i]
		meta := (*issue.IssueDetailLIST)[i] // Dereference to the MetaList, then access by slice index
		log.Println("  * ", issue.ID, issue.Subject)
		log.Println("  *  Meta :: ID", meta.ID)
		fmt.Printf("  * issue.FinishedDate = %s\n\n", issue.FinishedDate)
		fmt.Printf("  * meta.FinishedDate = %s\n\n", meta.FinishedDate)
	}

}
