/* Unceremoniously copied from Docker:
 * https://github.com/docker/docker/blob/master/pkg/namesgenerator/names-generator.go
 */
package namesgenerator

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomName(retry int) string {
	rand.Seed(time.Now().UnixNano())

	name := fmt.Sprintf("%s-%s", randomLeftSide(), randomRightSide())

	if retry > 0 {
		name = fmt.Sprintf("%s%d", name, rand.Intn(10))
	}

	return name
}

func Possibilities() int64 {
	return int64(len(leftSides) * len(rightSides))
}

func randomLeftSide() string {
	return leftSides[rand.Intn(len(leftSides))]
}

func randomRightSide() string {
	return rightSides[rand.Intn(len(rightSides))]
}
