# Test assignment

## Implementation details   
The API is designed to support multiple data storage types. Since the data must be persistent, the assumption that:

**_The app doesn't need to scale and there will always be a unique instance running at any given time_**

was made for simplicity sake due to the local nature of datastore that is implemented. The use of systems such as a database or a key/pair service would alleviate this restriction.

In this implementation, users and groups are kept in memory through a hash table (map).

The data stays in memory as long as the app runs and is only saved at shutdown in a local file after being serialized. When the app restarts, It reloads the data from the saved file.

The rationale for this implementation is that it doesn't require to install and run third party system (MySQL, redis etc...) to enable this app to work. This makes it more self contained and faster to get up and running which shouild make its evaluation easier.


## Installation & build
clone this repo

from the _planet_ folder, type
```text
> go build
```

## Running the app
The app serves client requests on port 80. 
To start the app, type
```text
> ./planet
```
The app will show these messages:

```text
SIGTERM and SIGINT signals initialized!
Waiting to shutdown
```
These messages mean that the SIGTERM and SIGINT signals have been captured to detect a shutdown and an appropriate action (save data to local file) will be taken when this happens.


## Testing the app
In one terminal, run the app. 
In a second terminal, go to the _api_ subfolder
type
```
> go test -v
``` 

Current tests only consider http status to determine whether a call succeeded or failed as expected. A deeper testing could be accomplished for endpoints returning a payload by also comparing their payload to what is expected. Due to time constraints to prepare scenarios that reflect a thorough testing, these features were left on the sideline.


## Introduction

This app implements a REST service used to store, fetch, and update user records. 
A user record can be represented in a JSON hash as follows: 

```json
{     
    "first_name": "Joe",     
    "last_name": "Smith",     
    "userid": "jsmith",     
    "groups": ["admins", "users"] 
}
```

This API helps maintains information about a set of users and the groups they are members of. HTTP verbs are used to implicitly indicate what the intent is, whether it is about _users_ or _groups_:

### GET
_GET_ is used to fetch data.

### POST
_POST_ is used to add new information.

### PUT
_PUT_ is used for data updates.

### DELETE
_DELETE_ is used to remove information.


An endpoint is characterized by its name and the http verb required for the call. Endpoints can be categorized in 2 classes: The **users** and **groups** classes. 

Endpoints of a given class use the same endpoint name, but require different http method so that they can be differentiated.


# API

## users class
These endpoints use **/users** as a verb.


- **GET /users/{userid}**     
Returns the matching user record or 404 if none exist.  

```text
GET /users/joes
```


- **POST /users**     
Creates a new user record. The body of the request should be a valid user record. 

```text
POST /users

{     
    "first_name": "Joe",     
    "last_name": "Smith",     
    "userid": "joes",     
    "groups": ["admins", "users"] 
}
```
POSTs to an existing user will return a 500 and the message _error user <userid> already exists_.



- **DELETE /users/{userid}**   
Deletes a user record.  

```text
DELETE /users/joes
```
Returns 404 if the user doesn't exist. 



- **PUT /users/{userid}**   
Updates an existing user record. The body of the request should be a valid user record. 

```text
PUT /users/joes

{     
    "first_name": "Joe",     
    "last_name": "Smith",     
    "userid": "joes",     
    "groups": ["admins", "users", "other"] 
}
```
PUTs to a non-existent user should return a 404.  



## groups class   
These endpoints use **/groups** as a verb.


- **GET /groups/{groupname}**   
Returns a JSON list of the members of that group. 

```text
GET /groups/admin
```

Returns a 404 if the group doesn't exist.

example of response:
```json
  ["userid1","userid2"]
```


- **POST /groups**   
Creates an empty group. 

```text
POST /groups

{     
  "group_name": "newgroup"
}

```

POSTs to an existing group will generate a 500. The body should contain a name parameter:


- **PUT /groups/{groupname}**   
Updates the membership list for the group. The body of the request should be a JSON list 
describing the group's members. 

```text
PUT /groups/admin

["userid1","userid2"]

```

Group members of the old list which are not in the new list will have the group removed from their _groups_ field. Similarly, the group name is added for members of the new list that did not belong to the old list. The body should contain an array of userid.

```json
  ["userid1","userid2"]
```

Returns a 404 if the group does not exist.


- **DELETE /groups/{groupname}**   
Deletes a group. 

```text
DELETE /groups/admin

```

Users member of the group _{groupname}_ get their _groups_ field updated to reflect the change.
