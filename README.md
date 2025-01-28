# WaterSportsRental

**WaterSportsRental** — веб-приложение для управления арендой водного инвентаря.

## Возможности
1. **Каталог товаров**: 
   - Фильтрация по заданным параметрам.
2. **Бронирование**: 
   - Выбор и бронирование инвентаря.
3. **Система отзывов**: 
   - Оставление комментариев о забронированных товарах.
4. **Личный кабинет**: 
   - Управление профилем и просмотр истории бронирований.
5. **Панель администратора**: 
   - Управление товарами, пользователями и резервными копиями базы данных.

## Технические особенности
1. **Фронтенд**:
   - Реализован на React с использованием TypeScript.
2. **Бэкенд**:
   - Написан на Golang с использованием роутера `chi`.
   - Организована микросервисная архитектура с использованием паттерна «репозиторий».
   - Авторизация реализована с помощью JWT-токенов.
3. **База данных**:
   - PostgreSQL, database/sql, lib/pq.
   - Для инициализации и удаления таблиц используются миграции.
   - Обработка запросов реаизована с помощью хранимых функций и процедур, триггеров, транзакций.
4. **Резервное копирование и восстановление базы данных**:
   - Возможность создания и восстановления резервных копий через серверное API.
   - Использование Docker для управления процессами `pg_dump` и `pg_restore`.

### Пример кода: Создание резервной копии

```go
func (h *Handler) createDump(w http.ResponseWriter, r *http.Request) {
    dumpCfg := dumpConfig.MustLoadDumpConfig()
    currentTime := time.Now()
    timeString := currentTime.Format("2006-01-02_15-04-05")
    filePath := filepath.Join(dumpCfg.Dir, fmt.Sprintf("%s_%s", dumpCfg.Prefix, timeString))

    cmd := exec.Command("docker", "exec", dumpCfg.ContainerName, "pg_dump", "-U", dumpCfg.Username, "-F", "c", dumpCfg.DbName)
    outputFile, err := os.Create(filePath)
    if err != nil {
        log.Fatal(err)
    }
    defer outputFile.Close()

    cmd.Stdout = outputFile

    if err = cmd.Run(); err != nil {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
    }
    err = h.services.Dump.InsertDump(r.Context(), filePath)
    if err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}
```
5. **Логирование**:
   - Поддерживаются уровни логирования: local, dev, prod.
6. **Graceful Shutdown**:
   - Реализован корректный процесс завершения работы сервера.
7. **Конфигурация**:
   - Конфиги для базы данных и сервера хранятся в отдельных файлах.
   - Приватные данные защищены в .env.

---

---
