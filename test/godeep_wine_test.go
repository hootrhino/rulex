package test

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)

// go test -timeout 30s -run ^Test_wine_demo github.com/hootrhino/rulex/test -v -count=1

func Test_wine_demo(t *testing.T) {

	rand.Seed(time.Now().UnixNano())

	data, err := load("./data/wine.data")
	if err != nil {
		panic(err)
	}

	for i := range data {
		deep.Standardize(data[i].Input)
		t.Log(data[i].Input[0], data[i].Response)
	}
	data.Shuffle()

	fmt.Printf("have %d entries\n", len(data))

	neural := deep.NewNeural(&deep.Config{
		Inputs:     len(data[0].Input),
		Layout:     []int{5, 3},
		Activation: deep.ActivationSigmoid,
		Mode:       deep.ModeMultiClass,
		Weight:     deep.NewNormal(1, 0),
		Bias:       true,
	})
	// t.Log(neural.String())
	//trainer := training.NewTrainer(training.NewSGD(0.005, 0.5, 1e-6, true), 50)
	//trainer := training.NewBatchTrainer(training.NewSGD(0.005, 0.1, 0, true), 50, 300, 16)
	//trainer := training.NewTrainer(training.NewAdam(0.1, 0, 0, 0), 50)
	trainer := training.NewBatchTrainer(training.NewAdam(0.1, 0, 0, 0), 50, len(data)/2, 12)
	//data, heldout := data.Split(0.5)
	trainer.Train(neural, data, data, 10000)
	// testData1 := []float64{13.48, 1.81, 2.41, 20.5, 100, 2.7, 2.98, .26, 1.86, 5.1, 1.04, 3.47, 920}
	// testData2 := []float64{12.37, 1.21, 2.56, 18.1, 98, 2.42, 2.65, .37, 2.08, 4.6, 1.19, 2.3, 678}
	testData3 := []float64{14.13, 4.1, 2.74, 24.5, 96, 2.05, .76, .56, 1.35, 9.2, .61, 1.6, 560}
	result1 := neural.Predict(testData3)
	result2 := neural.Predict(testData3)
	result3 := neural.Predict(testData3)
	t.Log(result1)
	t.Log(result2)
	t.Log(result3)
}
func one_hot(val float64) []float64 {
	// val 1,2,3
	// println("one_hot ==> ", classes, val)
	// res := make([]float64, classes)
	// res[int(val)-1] = 1
	if val == 1 {
		return []float64{0, 0, 1}
	}
	if val == 2 {
		return []float64{0, 1, 0}
	}
	if val == 3 {
		return []float64{1, 0, 0}
	}
	return []float64{0, 0, 0}

}

func load(path string) (training.Examples, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(bufio.NewReader(f))

	var examples training.Examples
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		examples = append(examples, toExample(record))
	}

	return examples, nil
}

func toExample(in []string) training.Example {
	// 1,14.23,1.71,2.43,15.6,127,2.8,3.06,.28,2.29,5.64,1.04,3.92,1065
	tag, err := strconv.ParseFloat(in[0], 64)
	if err != nil {
		panic(err)
	}
	tagEncoded := one_hot(tag)
	var params []float64
	for i := 1; i < len(in); i++ {
		param, err := strconv.ParseFloat(in[i], 64)
		if err != nil {
			panic(err)
		}
		params = append(params, param)
	}

	return training.Example{
		Input:    params,
		Response: tagEncoded,
	}
}

func Test_Prediction(t *testing.T) {
	rand.Seed(0)
	var data = []training.Example{
		{Input: []float64{2.7810836, 2.550537003}, Response: []float64{0}},
		{Input: []float64{1.465489372, 2.362125076}, Response: []float64{0}},
		{Input: []float64{3.396561688, 4.400293529}, Response: []float64{0}},
		{Input: []float64{1.38807019, 1.850220317}, Response: []float64{0}},
		{Input: []float64{3.06407232, 3.005305973}, Response: []float64{0}},
		{Input: []float64{7.627531214, 2.759262235}, Response: []float64{1}},
		{Input: []float64{5.332441248, 2.088626775}, Response: []float64{1}},
		{Input: []float64{6.922596716, 1.77106367}, Response: []float64{1}},
		{Input: []float64{8.675418651, -0.242068655}, Response: []float64{1}},
		{Input: []float64{7.673756466, 3.508563011}, Response: []float64{1}},
	}

	n := deep.NewNeural(&deep.Config{
		Inputs:     2,
		Layout:     []int{2, 2, 1},
		Activation: deep.ActivationSigmoid,
		Weight:     deep.NewUniform(0.5, 0),
		Bias:       true,
	})
	trainer := training.NewBatchTrainer(training.NewAdam(0.1, 0, 0, 0), 50, len(data)/2, 12)

	// trainer := NewTrainer(NewSGD(0.5, 0.1, 0, false), 0)

	trainer.Train(n, data, data, 5000)

	for _, d := range data {
		t.Log(n.Predict(d.Input))
	}
}
