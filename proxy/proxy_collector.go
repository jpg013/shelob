package proxy

import (
	"sync"

	"golang.org/x/net/html"
)

type proxyBuilder struct {
	root     *html.Node
	tds      []*html.Node
	styles   *html.Node
	siteName string
	data     map[string]interface{}
}

// Executor func takes a *proxyBuilder, modifies it and returns the pointer
type Executor func(*proxyBuilder) (*proxyBuilder, error)

// Pipeline interface represents a pipeline type with a Pipe and Merge method
type Pipeline interface {
	Pipe(fn Executor) Pipeline
	Merge() chan *proxyBuilder
}

// Collector implements pipeline interface with a list of executors
type Collector struct {
	executors []Executor
	dataCh    chan *proxyBuilder
}

// Merge chains the executors together and returns the output channel
func (c *Collector) Merge() chan *proxyBuilder {
	for i := 0; i < len(c.executors); i++ {
		fn := c.executors[i]
		c.dataCh = c.runExecutor(fn, c.dataCh)
	}

	return c.dataCh
}

// Pipe adds executor to the pipeline
func (c *Collector) Pipe(fn Executor) Pipeline {
	c.executors = append(c.executors, fn)
	return c
}

func (c *Collector) runExecutor(fn Executor, in <-chan *proxyBuilder) chan *proxyBuilder {
	out := make(chan *proxyBuilder)

	go func() {
		var wg sync.WaitGroup

		for p := range in {
			wg.Add(1)
			go func(p *proxyBuilder) {
				val, _ := fn(p)
				out <- val
				wg.Done()
			}(p)
		}

		go func() {
			defer close(out)
			wg.Wait()
		}()
	}()

	return out
}

// NewCollector factory returns a new Pipeline
func NewCollector(in chan *proxyBuilder) Pipeline {
	return &Collector{
		dataCh:    in,
		executors: make([]Executor, 0),
	}
}
