#!/bin/bash
# Script for me to easily update docker package
# ...I really need CI/CD
docker build https://git.a71.su/Andrew71/hibiscus.git -t git.a71.su/andrew71/hibiscus:latest
docker login git.a71.su
docker push git.a71.su/andrew71/hibiscus:latest