Kalandra
========

[![Build Status](https://travis-ci.org/kstrempel/kalandra.svg?branch=master)](https://travis-ci.org/kstrempel/kalandra)

## Idea 

The main idea of Kalandra is to provide a simple REST API for cassandra. 

## API

The REST interface is following the [JSONAPI](http://jsonapi.org/) specification. 

A query looks like the following

    GET /query
    {
        "data": {
            "keyspace": "dev",
            "attributes": {
                "query": "select * from emp;"
            }
        }
    }

The example returns the following answer. 

    {
        "data":[
            {
                "emp_dept":"eng",
                "emp_first":"fred",
                "emp_last":"smith",
                "empid":2
            }
        ],
        "Meta":{
            "query":"select * from emp;",
            "Time":34416848
        }
    }


## Build

Kalandra uses the great [gocql](https://github.com/gocql/gocql) implementation to access cassandra and mux for the web serving part.

    go get github.com/gocql/gocql
    go get github.com/gorilla/mux
    go get github.com/Sirupsen/logrus

    go build


## Testing

To run the unit tests a running cassandra server listening on localhost is needed.
Start the tests with:

    go test

License
-------

> Copyright (c) 2016 by Kai Strempel.
> Distributed under the Apache License
> license can be found in the LICENSE.txt file.
