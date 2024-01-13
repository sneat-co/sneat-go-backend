# Contactus module for Sneat.app

## [Models](./models) - database models

- [ContactDto](./models4contactus/contact_dto.go) - contact record data model
- [ContactusTeamDto](./models4contactus/contactus_team_dto.go) - contactus module team data record

## HTTP Endpoints

- [/v0/contactus/create_contact](./api4contactus/http_create_contact.go) - creates a new team contact
- [/v0/contactus/create_member](./api4contactus/http_create_member.go) - creates a new team member
