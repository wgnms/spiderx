package pandas
import (
	"github.com/go-gota/gota/series"
	"github.com/go-gota/gota/dataframe"
	"fmt"
)
func Test() {
	series.New([]string{"z", "y", "d", "e"}, series.String, "col")

	dataFrame := dataframe.New(
		series.New([]string{"a", "b", "c", "d", "e"}, series.String, "alphas"),
		series.New([]int{5, 4, 2, 3, 1}, series.Int, "numbers"),
		series.New([]string{"a1", "b2", "c3", "d4", "e5"}, series.String, "alnums"),
		series.New([]bool{true, false, true, true, false}, series.Bool, "state"),
	)
	fmt.Println(dataFrame)
}
