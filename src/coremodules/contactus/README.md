# Contactus module for Sneat.app

## [Models](./models) - database models

- [ContactDbo](./dbo4contactus/contact_dbo.go) - contact record data model
- [ContactusSpaceDbo](./dbo4contactus/contactus_space_dbo.go) - contactus module team data record

## HTTP Endpoints

- [/v0/contactus/create_contact](./api4contactus/http_create_contact.go) - creates a new team contact
- [/v0/contactus/create_member](./api4contactus/http_create_member.go) - creates a new team member
