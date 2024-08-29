package idgenerator

import "github.com/google/uuid"

type UuidIdGenerator struct{}

func NewUuidIdGenerator() *UuidIdGenerator {
	return &UuidIdGenerator{}
}

func (u *UuidIdGenerator) GenerateId() string {
	return uuid.NewString()
}
