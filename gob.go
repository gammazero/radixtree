package radixtree

import (
	"bytes"
	"encoding/gob"
)

func (t *Tree[T]) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	var err error

	if err = encoder.Encode(&t.root); err != nil {
		return nil, err
	}
	if err = encoder.Encode(&t.size); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (t *Tree[T]) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	var err error

	if err = decoder.Decode(&t.root); err != nil {
		return err
	}
	return decoder.Decode(&t.size)
}

func (node *radixNode[T]) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(node.prefix)
	if err != nil {
		return nil, err
	}
	if err = encoder.Encode(node.edges); err != nil {
		return nil, err
	}
	if node.leaf != nil {
		if err = encoder.Encode(node.leaf); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (node *radixNode[T]) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	err := decoder.Decode(&node.prefix)
	if err != nil {
		return err
	}
	if err = decoder.Decode(&node.edges); err != nil {
		return err
	}
	if buf.Len() != 0 { // if leaf is not nil
		return decoder.Decode(&node.leaf)
	}
	return nil
}

func (kv *Item[T]) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(kv.key)
	if err != nil {
		return nil, err
	}
	if err = encoder.Encode(kv.value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (kv *Item[T]) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	err := decoder.Decode(&kv.key)
	if err != nil {
		return err
	}
	return decoder.Decode(&kv.value)
}

func (e *edge[T]) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(e.radix)
	if err != nil {
		return nil, err
	}
	if err = encoder.Encode(e.node); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *edge[T]) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	err := decoder.Decode(&e.radix)
	if err != nil {
		return err
	}
	return decoder.Decode(&e.node)
}
