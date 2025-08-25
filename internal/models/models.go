package models

import (
	"time"

	"github.com/google/uuid"
)

// надо объединить это и кол-во изображений, а то когда надо, то несколько запросов придется делать ради одного числа
// это осознанная денормализация

// да, можно подумать о денормализации, но кол-во изображений мне нужно только при добавлении новых,
// поэтому для того, чтобы несколько раз не идти в БД, при добавлении нескольких фотографий не сразу
// использую объединенную структуру в кэше со всей информацией

// все предыдущие записи недействительны - информация о максимальном и текущем кол-ве изображений
// и так хранится в сервисе изображений, всё равно делать запрос, поэтому в сервисном слое тут
// якобы определяю по пользователю, сколько ему можно загружать картинок, и просто отправляю эту информацию
// обратно в шлюз, чтобы использовать её в следующем запросе к сервису картинок
type Product struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	Price       int
	CreatedAt   time.Time
}

type NewProduct struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	Price       int
}

// так эта информация уже есть в сервисе изображений
/*
type ProductImage struct {
	ProductID uuid.UUID
	MaxCount  int
}

type FullProductInfo struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	Price       int
	CreatedAt   time.Time
	MaxCount    int
}
*/
