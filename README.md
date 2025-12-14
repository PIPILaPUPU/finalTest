# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

#Ниже описание работы
Было разработан backend для приложения по созданию онлайн дневника
Из задач со звездочкой были выполнены:
-Реализация чтения порта из переменной окружения TODO_PORT
-Нахождение задачи через строку. Например если в поиск вбить бассейн, то выведутся все задачи со словом бассейн

Флаги:
Для запуска стоит установить значение переменной окружения TODO_PORT=7540, а для TODO_PASSWORD=12345. Именно такие значения я использовал при тестировании

Параметры settings:
var Port = GetPort()
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwd2RfaGFzaCI6IjU5OTQ0NzFhYmIwMTExMmFmY2MxODE1OWY2Y2M3NGI0ZjUxMWI5OTgwNmRhNTliM2NhZjVhOWMxNzNjYWNmYzUiLCJleHAiOjE3NjU3NTc1MzgsImlhdCI6MTc2NTcyODczOH0.u5kKrEPEZaa2v13aVSGMJRsoXxEb7y1ScrQJHOj1OUQ`
данный токен применялся при пароле TODO_PASSWORD=12345

func GetPort() int {
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			return port
		}
	}

	return 7540
}

Запуск:
Есть два способа запуска, с паролем и без него. 
С паролем go run main.go
Без пароля TODO_PASSWORD="" go run main.go

При запуске с докера использовались следующие теги
docker run -p 7540:7540 todo-app:latest
