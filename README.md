# App-Metadata

- [App-Metadata](#app-metadata)
  - [Description](#description)
  - [Getting started](#getting-started)
  - [External dependencies](#external-dependencies)
  - [Packages](#packages)
    - [handlers](#handlers)
    - [metadata](#metadata)
    - [repository](#repository)

## Description

App metadata is a REST-API to store a metadata in yaml format

## Getting started

- It requires go version >= 1.11 to enable go modules
- Clone it from <https://github.com/elumbantoruan/app-metadata.git> (https) or git@github.com:elumbantoruan/app-metadata.git (ssh)
- To build, run go build ./...
- To run unit test, run go test ./...
- To run the app, execute go run main.go.  This will enable the endpoint at localhost:5000/app-metadata
- Example of POST operation returns 201, and the created payload

``` text
curl -d '
> title: Valid App 1
> version: 1.0.1
> maintainers:
> - name: First Maintainer App1
>   email: firstmaintainer@hotmail.com
> - name: Second Maintainer App1
>   email: secondmaitainer@gmail.com
> company: pellucid Computing
> website: http://pellucidcomputing.com
> source: https://github.com/company/app-metadata
> license: Apache-2.0
> description: |-
>   ### Interesting title
>   Some application content' -i http://localhost:5000/app-metadata
```

- Example of POST operation returns 400 because of invalid payload

``` text
curl -d '
> title: Valid App 1
> maintainers:
> - name: First Maintainer App1
>   email: firstmaintainer@hotmail.com
> - name: Second Maintainer App1
>   email: secondmaitainer@gmail.com
> company: pellucid Computing
> website: http://pellucidcomputing.com
> source: https://github.com/company/app-metadata
> license: Apache-2.0
> description: |-
>   ### Interesting title
>   Some application content' -i http://localhost:5000/app-metadata
```

## External dependencies

Go modules (go.mod) will simply download all dependencies after running the build command

- [Gorilla mux](https://github.com/gorilla/mux)  It's URL router and dispatcher
- [UUID](https://github.com/google/uuid) It's a UUID to generate a unique id
- [Testify](https://github.com/stretchr/testify) Tools for unit test such as assert, suite, and mock]
- [yaml](gopkg.in/yaml.v2) YAML support for the Go language

## Packages

### handlers

MetadataHandler contains a set of methods to perform POST, PUT, GET, and DELETE for REST API operations.

``` text
POST   /app-metadata
PUT    /app-metadata/{appID}
GET    /app-metadata
GET    /app-metadata/{appID}
DELETE /app-metadata/{appID}
```

It has a dependency on repository interface to perform a create, get, update, and delete repository actions.

It also contains a unit test for all REST API operations

### metadata

ApplicationMetadata is a payload used in the application which is marshalled into yaml format

### repository

Repository contains a MetadataRepository interface and InMemoryMetadataRepository type.  The intent of the interface is to allow flexibility for swapping different repository mechanisms.  It also enables a mock up repository to be used in unit test

``` go
type MetadataRepository interface {
    Create(appID string, data *metadata.ApplicationMetadata) error
    Update(appID string, data *metadata.ApplicationMetadata) error
    Get(appID string) (*metadata.ApplicationMetadata, error)
    GetAll() ([]metadata.ApplicationMetadata, error)
    Delete(appID string) error
}
```

InMemoryMetadataRepository is a concrete implementation of MetadataRepository interface
