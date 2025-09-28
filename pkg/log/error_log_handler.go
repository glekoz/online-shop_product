package log

import (
	"context"
	"errors"
)

type ErrorLogData struct {
	LD  LogData
	Err error
}

// а нужен ли мне тут указатель? придется хранить результат функции в куче
// состояние не изменяется к тому же
// func (e *ErrorLogData) Error() string {
func (e ErrorLogData) Error() string {
	return e.Err.Error()
}

func WrapError(ctx context.Context, err error) error {
	ld := LogData{}
	if ldt, ok := ctx.Value(LogDataKey).(LogData); ok {
		ld = ldt
	}
	return ErrorLogData{LD: ld, Err: err}
}

func ErrorContext(ctx context.Context, err error) context.Context {
	var errt ErrorLogData
	if errors.As(err, &errt) {
		return context.WithValue(ctx, LogDataKey, errt.LD)
	}
	return ctx
}

/*
func Handler(ctx context.Context, userID int) {
	ctx = WithUserID(ctx, userID)
	slog.InfoContext(ctx, "Handler started")
	phone, err := GetPhoneByID(ctx, userID)
	if err != nil {
		slog.ErrorContext(ErrorContext(ctx, err), "Error: "+err.Error())
		return
	}
	ctx = WithPhone(ctx, phone)
	err = SendSMS(ctx, phone)
	if err != nil {
		slog.ErrorContext(ErrorContext(ctx, err), "Error: "+err.Error())
		return
	}
	slog.InfoContext(ctx, "Handler done")
}
*/
