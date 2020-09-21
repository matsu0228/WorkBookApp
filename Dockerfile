FROM golang:latest

RUN mkdir /GolandProjects/WorkBookApp

WORKDIR /GolandProjects/WorkBookApp

ADD . /GolandProjects/WorkBookApp
