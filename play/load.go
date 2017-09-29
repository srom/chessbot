package play

import (
	"io/ioutil"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type Model struct {
	graph   *tf.Graph
	session *tf.Session
	x       tf.Output
	y       tf.Output
}

func (m *Model) Evaluate(input BoardInput) (int64, error) {
	inputTensor, err := tf.NewTensor(input)
	if err != nil {
		return 255, err
	}
	tensors, err := m.session.Run(
		map[tf.Output]*tf.Tensor{
		    m.x: inputTensor,
		},
		[]tf.Output{m.y},
		nil,
	)
	if err != nil {
		return 255, err
	}
	output := (tensors[0].Value().([]int64))[0]
	return output, nil
}

func (m *Model) Close() error {
	return m.session.Close()
}

func LoadModel(modelPath string) (*Model, error) {
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

	model := Model{
		graph: graph,
		session: session,
		x: graph.Operation("chessbot/prediction_inputs").Output(0),
		y: graph.Operation("chessbot/dnn/head/predictions/class_ids").Output(0),
	}
    	return &model, nil
}
