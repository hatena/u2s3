#!/bin/bash

exec docker run -p 9000:9000 --name log2s3_minio_test \
       -e "MINIO_ACCESS_KEY=ACCESS_KEY" \
       -e "MINIO_SECRET_KEY=SECRET_KEY" \
       minio/minio server /export
