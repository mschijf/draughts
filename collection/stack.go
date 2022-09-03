package collection

type Stack[T any] struct {
	stack []*T
}

func (thisStack *Stack[T]) Push(element *T) {
	thisStack.stack = append(thisStack.stack, element)
}

func (thisStack *Stack[T]) Pop() *T {
	n := len(thisStack.stack) - 1 // Top element
	element := thisStack.stack[n]
	thisStack.stack = thisStack.stack[:n]
	return element
}

func (thisStack *Stack[T]) IsEmpty() bool {
	return len(thisStack.stack) == 0
}

func (thisStack *Stack[T]) Top() *T {
	return thisStack.FromTop(0)
}

func (thisStack *Stack[T]) FromTop(i int) *T {
	return thisStack.stack[len(thisStack.stack)-i-1]
}

func (thisStack *Stack[T]) Size() int {
	return len(thisStack.stack)
}
