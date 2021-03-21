# AstraTutor

This repo serves as the monorepo for the AstraTutor application for Team 4 of the CS3305 Team Software Project module

## Instructions

We have 2 test accounts for logging in
```
Tutor email:        tutor@grindsapp.localhost
Tutor password:     grindsapp

Student email:      student@grindsapp.localhost
Student password:   grindsapp
```

To look at pg admin use the following credentials
```
Email:              grindsapp@grindsapp.localhost
Password:           grindsapp

Second password:    grindsapp
```

### Use website
Visit [astratutor.com](https://astratutor.com/)

### Install Locally
1. Clone the repo
1. Run `docker-compose up` and visit:
    * `https://localhost:8080` for UI
    * `https://localhost:8081` for API
    * `https://localhost:8082` for pgAdmin

This will take some time to start as it needs to seed the database

## Folder Structure
```
.
|-- api                 Contains GO Backend
|   |-- cmd             Contains script used to start webserver
|   |-- pkg             Contains backend code
|   |   |-- database    Contains code to open database connection
|   |   |-- routes      Contains code to define api endpoints
|   |   `-- services    Contains services code
|   `-- seed            Contains data used for seeding
|-- documents           Contains various documents for project
|-- pgadmin4            Contains config for pgadmin
`-- ui                  Contains React frontend
    |-- public          Contains Static files
    `-- src             Contains frontend components
        |-- api         Contains typescript for interfacing with backend
        |-- components  Contains react components
        |-- views       Contains pages
        `-- webrtc      Contains typescript for webrtc connections
```

## Documentation
https://docs.astratutor.com/v1.pdf
