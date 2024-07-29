package slurm
//
// import (
// 	"io/ioutil"
// 	"os"
// 	"testing"
// )
//
// func TestParseQueueMetrics(t *testing.T) {
// 	// Read the input data from a file
// 	file, err := os.Open("test_data/squeue.txt")
// 	if err != nil {
// 		t.Fatalf("Can not open test data: %v", err)
// 	}
// 	data, err := ioutil.ReadAll(file)
// 	t.Logf("%+v", ParseQueueMetrics(data))
// }
//
// func TestQueueGetMetrics(t *testing.T) {
// 	t.Logf("%+v", QueueGetMetrics())
// }
