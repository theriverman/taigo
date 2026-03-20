package taigo

// TaigaBaseObject represents the following Taiga object types:
/*
* Epic
* User Story
* Task
* Issue

These Taiga objects have the following must-have fields in common:
* ID
* Ref
* Version
* Subject
* Project
*/
type TaigaBaseObject interface {
	GetID() int
	GetRef() int
	GetVersion() int
	GetSubject() string
	GetProject() int
}

// AttachmentTarget is the minimum contract required to upload attachments.
type AttachmentTarget interface {
	GetID() int
	GetProject() int
}
