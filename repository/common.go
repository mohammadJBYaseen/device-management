package repository

import (
	"device-management/model"
	"fmt"
	"strings"
)

func Order(sort model.Sort) string {
    return fmt.Sprintf("%s %s", sort.SortBy, strings.ToLower(sort.Direction));
}
