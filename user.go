package main

import "errors"

type Service struct{}

func (s *Service) GetUserName(id int) (string, error) {
    if id == 0 {
        return "", errors.New("invalid user id")
    }
    return "John", nil
}