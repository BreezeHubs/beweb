package beweb

// Middleware 函数式责任链模式（洋葱模式）
type Middleware func(next HandleFunc) HandleFunc
