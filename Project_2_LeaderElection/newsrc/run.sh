#!/bin/bash
./newsrc :8080       :8081 :8082 :8083 :8084 &
./newsrc :8081 :8080       :8082 :8083 :8084 &
./newsrc :8082 :8080 :8081       :8083 :8084 &
./newsrc :8083 :8081 :8081 :8082       :8084 &
./newsrc :8084 :8081 :8081 :8082 :8083       &
