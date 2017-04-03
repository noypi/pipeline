/*
	See also, https://github.com/noypi/fp
*/
package util

type Pipe func(in, out chan interface{})

type Pipeline []Pipe

func (pipes Pipeline) Run(in chan interface{}, size int) (out interface{}) {
	if 0 >= size {
		size = 1
	}
	chs := append([]chan interface{}{in}, make([]chan interface{}, len(pipes))...)
	for i := 1; i < len(chs); i++ {
		chs[i] = make(chan interface{}, size)
	}

	for i, j := 0, 1; j < len(chs); i, j = i+1, j+1 {
		go pipes[i](chs[i], chs[j])
	}

	out = <-chs[len(chs)-1]

	for _, ch := range chs[1:] {
		close(ch)
	}

	return

}

func PipeForwarder(in, out chan interface{}) {
	for v := range in {
		out <- v
	}
}

func (pipes Pipeline) RunAsyncOut(in, out chan interface{}, size int) (cleanup func()) {
	if 0 >= size {
		size = 1
	}

	// create temp channels
	chs := append([]chan interface{}{in}, make([]chan interface{}, len(pipes))...)
	for i := 1; i < len(chs); i++ {
		chs[i] = make(chan interface{}, size)
	}

	for i, j := 0, 1; j < len(chs); i, j = i+1, j+1 {
		go pipes[i](chs[i], chs[j])
	}

	go PipeForwarder(chs[len(chs)-1], out)

	return func() {
		//exclude in
		for _, ch := range chs[1:] {
			close(ch)
		}
	}

}
