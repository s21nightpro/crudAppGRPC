#!/bin/bash

lsof -ti:50057 | xargs kill -9
go run server.go