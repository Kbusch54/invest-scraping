package httputil

type Response[T any] struct {
	Msg      string `json:"msg"`
	HttpCode int    `json:"http_code"`
	Success  bool   `json:"success"`
	Data     T      `json:"data"`
}

func OK[T any](msg string, data T) Response[T] {
	return Response[T]{
		Msg:      msg,
		HttpCode: 200,
		Success:  true,
		Data:     data,
	}
}

func Created[T any](msg string, data T) Response[T] {
	return Response[T]{
		Msg:      msg,
		HttpCode: 201,
		Success:  true,
		Data:     data,
	}
}

func NotFound[T any](msg string, data T) Response[T] {
	return Response[T]{
		Msg:      msg,
		HttpCode: 404,
		Success:  false,
		Data:     data,
	}
}

func InternalServerError[T any](data T) Response[T] {
	return Response[T]{
		Msg:      "Internal server error :(",
		HttpCode: 500,
		Success:  false,
		Data:     data,
	}
}

func BadRequestError[T any](data T) Response[T] {
	return Response[T]{
		Msg:      "Bad Request :(",
		HttpCode: 400,
		Success:  false,
		Data:     data,
	}
}

func Conflict[T any](data T) Response[T] {
	return Response[T]{
		Msg:      "Conflict",
		HttpCode: 409,
		Success:  false,
		Data:     data,
	}
}

func Unauthorized[T any](data T) Response[T] {
	return Response[T]{
		Msg:      "Unauthorized :x",
		HttpCode: 401,
		Success:  false,
		Data:     data,
	}
}
