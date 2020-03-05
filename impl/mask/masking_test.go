package mask

import (
	"fmt"
	"testing"
)

type Data struct {
	Masking

	Head bool
	Body bool
}

func TestMasking(t *testing.T) {
	mask := byte(0)
	data := Data{
		Head: true,
	}

	data.Set(&mask, 0x01, data.Head)
	data.Set(&mask, 0x02, data.Body)

	fmt.Println(data.Has(mask, 0x01))
	fmt.Println(data.Has(mask, 0x02))
}
