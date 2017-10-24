package play

import (
	"io/ioutil"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type ChessBotModel struct {
	graph   *tf.Graph
	session *tf.Session
	x       tf.Output
	f       tf.Output
}

func (m *ChessBotModel) Evaluate(inputs [][]float32) ([]float32, error) {
	inputTensor, err := tf.NewTensor(inputs)
	if err != nil {
		return []float32{}, err
	}
	tensors, err := m.session.Run(
		map[tf.Output]*tf.Tensor{
		    m.x: inputTensor,
		},
		[]tf.Output{m.f},
		nil,
	)
	if err != nil {
		return []float32{}, err
	}
	output := []float32{}
	for _, tensor := range tensors[0].Value().([][]float32) {
		val := tensor[0]
		output = append(output, val)
	}
	return output, nil
}

func (m *ChessBotModel) Close() error {
	return m.session.Close()
}

func LoadModel(modelPath string) (*ChessBotModel, error) {
	modelBytes, err := ioutil.ReadFile(modelPath)
	if err != nil {
		return nil, err
	}

	graph := tf.NewGraph()
    	if err = graph.Import(modelBytes, "chessbot"); err != nil {
        	return nil, err
    	}

	session, err := tf.NewSession(graph, nil)
    	if err != nil {
        	return nil, err
    	}

	model := ChessBotModel{
		graph: graph,
		session: session,
		x: graph.Operation("chessbot/input/X").Output(0),
		f: graph.Operation("chessbot/f_p/output/BiasAdd").Output(0),
	}
    	return &model, nil
}
