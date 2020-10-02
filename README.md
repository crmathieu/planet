# PLANET - Code test assignment
 
Implement a REST service using a golang web framework that can be used to store, 
fetch, and update user records. A user record can be represented in a JSON hash 
like so: 

```json
{     
    "first_name": "Joe",     
    "last_name": "Smith",     
    "userid": "jsmith",     
    "groups": ["admins", "users"] 
}
```

This API uses http verbs to indicate their implicite behavior:

### GET
_GET_ is used to fetch information about a user or a group

### POST
_POST_ is used to add information such as new user or new group

### PUT
_PUT_ is used for information updates

### DELETE
_DELETE_ is used to remove a user or a group

An endpoint is characterized by its name and the http verb required for the call. 

There are 2 classes of endpoints: The **users** class, and the **groups** class. 

## API
The service should provide the following endpoints and semantics: 

### users class

**GET /users/<userid>**     
Returns the matching user record or 404 if none exist.  


**POST /users**     
Creates a new user record. The body of the request should be a valid user record. POSTs to an 
existing user will return a 500 and the message _error user <userid> already exists_.

**DELETE /users/<userid>**
Deletes a user record. Returns 404 if the user doesn't exist.  

**PUT /users/<userid>**
Updates an existing user record. The body of the request should be a valid user record. PUTs to 
a non-existent user should return a 404.  


### groups class

**GET /groups/<groupname>**
Returns a JSON list of user ids containing the members of that group. Returns a 404 
if the group doesn't exist.

example of response:
```json
  ["userid1","userid2"]
```

**POST /groups**   
Creates an empty group. POSTs to an existing group should be treated as errors and flagged 
with the appropriate HTTP status code. The body should contain a name parameter:

```json
{
  "group_name": "newgroup"
}
``` 

**PUT /groups/<groupname>**
Updates the membership list for the group. The body of the request should be a JSON list 
describing the group's members. The group name is removed for members of the old list but not in 
the new one. Similarly, the group name is added for members of the new group that were not in the old list.
The body should contain an array of userid.

```json
  ["userid1","userid2"]
```

**DELETE /groups/<groupname>**
Deletes a group. Users member of the group get their _groups_ information updated to reflect the change.


## Implementation
The API is design to allow for multiple type of repositories without requiring to change the unit test code. For simplicity and since the data must be persistent, I made the following assumptions: 

- The app doesn't need to scale and there will always be a unique instance running at any given time.
- I keep data in-memory and save it at shutdown in a local file. When the app restarts, It reloads the data from the save file.


Implementation Notes: 
1. Any design decisions not specified herein are fair game. Completed projects will be evaluated on how 
closely they follow the spec, their design, and cleanliness of implementation. 

2. The data stored should be durable. We love SQL at Planet but are happy for alternative methods 
of long term storage. 

3. Completed projects must include a README with enough instructions for evaluators to build and run 
the code. Bonus points for builds which require minimal manual steps. 

4. Remember this project should take a maximum of 8 hours to complete. Do not get hung up on scaling 
or persistence issues. This is a project used to evaluate your design and implementation skills only.

5. Please include any unit or integration tests used to verify correctness.



////////////////////////////////////////////////






Semaphores are tools that help manage concurrent accesses to a common resource (the resource is usually one or more data structures with a set of indexes or pointers used to manipulate the data). In order to achieve this, a semaphore has a counter that indicates the level of availability the resource has at a given time.

In the context of Golang, a semaphore will be represented as a data structure using a channel as a mean to provide goroutines synchronization. The channel has a dimension (its capacity) that corresponds to the dimension of the sharable resource, and an initial count corresponding to its initial availability.

### Using gosem in your project
```go
import (
  sem "github.com/crmathieu/gosem/pkg/semaphore"
)
```

Imagine we want to share a buffer of 512 integers between 1 producer and 1 consumer. The producer wants to write to the buffer and the consumer wants to read from it.

```go
var buffer [512]int
var tail, head int = 0, 0
```

In order to protect this buffer and synchronize read and write operations, we are going to need 2 semaphores: one for reads and one for writes. The read semaphore is used to find out if there is anything to read from the buffer. The write semaphore is used to find out if there is any space available in the buffer so that we can write into it.

```go
readsem =  sem.Createsem("readsem",  512, 0)
writesem = sem.Createsem("writesem", 512, 512)
```

Note that both readsem and writesem have the same dimension: 512, but readsem has an initial value of 0 (because initially there is nothing to read) and writesem has an initial value of 512 (because initially the whole buffer is available).

you may also use the _CreateReadSemaphore_ or _CreateWriteSemaphore_ that will abstract the initial value given to the semaphore:

```go
readsem =  sem.CreateReadSemaphore("readsem",   512)
writesem = sem.CreateWriteSemaphore("writesem", 512)
```

The code of the producer looks like this:

```go
func producer() {
  i := 0
  for {
    writesem.Wait()
    buffer[head] = i
    i = (i + 1) % 4096
    head = (head+1) % 512
    readsem.Signal()
  }
}
```
In its loop, the **producer** first makes sure there is available space in the buffer by calling _writesem.Wait()_. This call will return immediately if space is available but will block if the buffer is full. In the latter case, the call will return only after the **consumer** goroutine reads an entry from the buffer and performs a _writesem.Signal()_ call to indicate that one entry is now available.

Similarly, once a value was written in the buffer, the **producer** calls _readsem.Signal()_ to indicate that one entry is available for consumption.

The code of the consumer looks like that:

```go
func consumer() {
  for {
    readsem.Wait()
    item = buffer[tail]
    tail = (tail+1) % 512
    writesem.Signal()
    fmt.Println(item)
  }
}
```

In its loop, the **consumer** first makes sure there is something to read from the buffer by calling _readsem.Wait()_. This call will return immediately if data is available but will block if the buffer is empty. In the latter case, the call will return only after the **producer** goroutine writes an entry to the buffer and performs a _readsem.Signal()_ call to indicate that one entry is ready to be consumed.

Similarly, once a value has been read from the buffer, the **consumer** calls _writesem.Signal()_ to indicate that space is available for production.


## Multiple consumers and producers
If we want to synchronize several consumers and producers accessing the same buffer, the code for both consumers and producers needs to handle concurrent access to the <b>head</b> and <b>tail</b> buffer indexes, because their value can be updated by multiple goroutines simultaneously (which was not the case in the previous example).

In order to do that, goroutines will need to have an exclusive access to these indexes when they update them. This is accomplished with the use of <b>mutex semaphores</b>. A mutex semaphore is like a normal semaphore with a capacity and an initial count of 1:

```go
mutex = sem.Createmutex("mymutex")
```
We are going to need a mutex to protect the <b>head</b> index used by multiple producers and another mutex to protect the <b>tail</b> index used by multiple consumers:

```go
headmutex = sem.Createmutex("head-mutex")
tailmutex = sem.Createmutex("tail-mutex")
```

The producer code becomes:

```go
func producer() {
  i := 0
  for {
    writesem.Wait()
    headmutex.Enter()
    buffer[head] = i
    i = (i + 1) % 4096    
    head = (head+1) % 512
    headmutex.Leave()
    readsem.Signal()
  }
}
```

and the consumer code becomes:

```go
func consumer() {
  for {
    readsem.Wait()
    tailmutex.Enter()
    item = buffer[tail]
    tail = (tail+1) % 512
    tailmutex.Leave()
    writesem.Signal()
    fmt.Println(item)
  }
}
```

### Semaphore API

First, import the gosem package:

```go
import (
  sem "github.com/crmathieu/gosem/pkg/semaphore"
)
```

#### Variable declaration
To declare a semaphore or a mutex:
```go
var mysem *semaphore.Sem
```
-or-
```go
var mymutex *semaphore.Mutex
```

#### Createsem: creates a counter semaphore
_func Createsem(name string, capacity int, initialcount int) *Sem_

To create a semaphore with a capacity of 64, and an initial count of 0:
```go
mysem := sem.Createsem("mySemaphore", 64, 0)
```
-or- to create a semaphore with a capacity of 64, and an initial count of 64:
```go
mysem := sem.Createsem("mySemaphore", 64, 64)
```

#### Createmutex: creates a mutex  
_func Createmutex(name string) *Mutex_
```go
mymutex := sem.Createmutex("myMutex")
```

Following a semaphore creation, there are a certain number of methods available to manipulate semaphores:

#### Reset
_func (s *Sem) Reset()_
```go
mysem.Reset()
```
-or- for a mutex

_func (m *Mutex) Reset()_
```go
mymutex.Reset()
```

This will flush the semaphore internal channel and resets its counter to its original value.

#### Signal -or- V (-or- Leave)
_func (s *Sem) Signal()_
```go
mysem.Signal()
```
-or-

_func (s *Sem) V()_
```go
mysem.V()
```
-or- for a mutex

_func (m *Mutex) Leave()_
```go
mymutex.Leave()
```

<b>Signal</b> and <b>V</b> accomplish the same thing which is to increase by 1 the level of availability of the resource. <b>Leave</b> is identical but reserved for <i>mutex</i>.

#### Wait -or- P (-or- Enter)
_func (s *Sem) Wait()_
```go
mysem.Wait()
```
-or-

_func (s *Sem) P()_
```go
mysem.P()
```
-or- for a mutex

_func (m *Mutex) Enter()_
```go
mymutex.Enter()
```
<b>Wait</b> and <b>P</b> accomplish the same thing which is to decrease by 1 the level of availability of the resource. <b>Enter</b> is identical but reserved for <i>mutex</i>. When the semaphore counter reaches 0, the resource is no longer available, until a Signal (-or- a V) call is made by another goroutine.

#### Notes:
- The <b>P</b> / <b>V</b> notation comes from <b>Edsger Dijkstra</b>, who introduced the concept of semaphores in 1963. The letters are from the Dutch words <b>Probeer</b> (try) and <b>Verhoog</b> (increment).

- The terms <b>Enter</b> and <b>Leave</b> for a mutex refer to ```Entering``` and ```Leaving``` critical sections in your code. A <b>critical section</b> is a region in your code that can be executed only by one goroutine at a time. Typically, you will need to define a critical section everytime you need to access a resource that can potentially be modified by multiple goroutines. Once in the critical section, a goroutine is guaranteed to have exclusive access to the shared resource.

