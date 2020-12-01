package service

import (
	"Go-000/Week02/dao"

	"github.com/pkg/errors"
)

var ServiceErr = errors.New("Service not found")

// age limit
const AgeLimit = 30

type user struct {
	name string
	age  int
}

func Biz() error {
	u := &user{age: 30, name: "Li Ming"}
	err := dao.Dao()
	// if err != nil {
	//	return err
	// }

	if u.age <= AgeLimit {
		return errors.Wrap(err, "The age too old.")

	}

	return nil
}
