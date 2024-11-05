package logic

import (
	"fmt"
	"library-study/app/model"
	"testing"
)

func TestName(t *testing.T) {

	fmt.Println(model.GetPaginatedBooksData(1, 10))
}
