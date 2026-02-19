package main

import "testing"

func TestGetUserName_Success(t *testing.T) {
    svc := &Service{}

    name, err := svc.GetUserName(1)

    AssertNoError(t, err)
    AssertEqual(t, "John", name)
}

func TestGetUserName_InvalidID(t *testing.T) {
    svc := &Service{}

    name, err := svc.GetUserName(0)

    AssertEqual(t, "", name)
    AssertError(t, err, "invalid user id")
}