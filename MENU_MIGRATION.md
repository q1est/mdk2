# 🍽️ Миграция меню с фронта на бэк

## 📋 Краткая сводка

Меню переносится из хардкода HTML в динамическую загрузку с API. Все данные теперь хранятся в БД PostgreSQL и загружаются при открытии страницы.

## 📝 Что было изменено

### Frontend (JavaScript)

✅ **Новый файл**: `js/menuLoader.js`
- Загружает меню с endpoint'а `/api/menu`
- Динамически рендерит элементы меню
- Работает с корзиной и модалями

✅ **Обновлено**: `js/main.js`
- Добавлен импорт `menuLoader.js`

✅ **Обновлено**: `menu.html`
- Удалено хардкодированное меню
- Оставлена только структура (nav, main)

### Backend (Go)

✅ **Исправлено**: `back/db/db.go`
- Добавлен правильный SQL для создания `menu_items` таблицы

✅ **Исправлено**: `back/models/models.go`
- Исправлена опечатка `Descripition` → `Description`

✅ **Исправлено**: `back/db/postgres.go`
- Используется правильное имя поля `Description`

✅ **Существует**: `back/handle/handle.go`
- Endpoint `/api/menu` уже готов и возвращает JSON

## 🚀 Как это запустить

### Шаг 1: ~~Обновить URL API~~ ✅ ГОТОВО

Используется относительный путь `/api/menu` (как в order.js) — работает везде без изменений!

### Шаг 2: Убедиться что БД настроена

1. Создай/убедись что PostgreSQL запущена
2. Установи переменную окружения `DATABASE_URL`:
   ```bash
   export DATABASE_URL="postgres://user:password@localhost:5432/restaurant"
   ```

3. Запусти бэк-сервер (таблицы создадутся автоматически):
   ```bash
   go run back/main.go
   ```

### Шаг 3: Заполнить меню (первый раз)

Используй SQL скрипт `insert_menu.sql`:
```bash
psql $DATABASE_URL < insert_menu.sql
```

Или выполни в консоли БД:
```sql
INSERT INTO menu_items (name, description, price, category, image_url, available) VALUES
('Оливье', 'Оливье в жестяной банке...', 350, 'snacks', 'Оливье.jpg', true),
...
```

### Шаг 4: Проверить CORS

Если фронт на другом домене, убедись что в `handle.go` установлена правильная origin:

```go
// Локально
w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

// Production
w.Header().Set("Access-Control-Allow-Origin", "https://q1est.github.io")
```

### Шаг 5: Тестирование

1. Запусти бэк
2. Открой `menu.html` в браузере (можешь использовать Live Server)
3. Должно появиться меню, загруженное с API
4. Проверь консоль (F12) на ошибки
5. Добавь товар в корзину — должно работать

## 📊 Структура данных

### API ответ

```json
[
  {
    "id": 1,
    "name": "Оливье",
    "description": "Оливье в жестяной банке...",
    "price": 350,
    "category": "snacks",
    "image_url": "Оливье.jpg"
  },
  ...
]
```

### Categories

Категории (case-insensitive):
- `snacks` — Закуски
- `main` — Основные блюда
- `desserts` — Десерты
- `drinks` — Напитки

### Menu Items таблица

```sql
CREATE TABLE menu_items (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  price INT NOT NULL,
  category TEXT NOT NULL,
  image_url TEXT NOT NULL,
  available BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT NOW()
);
```

## 🐛 Возможные проблемы и решения

### ❌ Меню не загружается

**Проверь:**
1. Запущен ли бэк-сервер?
   ```bash
   curl http://localhost:7860/api/menu
   ```

2. Правильный ли URL в `menuLoader.js`?

3. CORS ошибка в консоли? Проверь `handle.go` Origin

4. БД подключена? Проверь `DATABASE_URL`

### ❌ Ошибка `desctripton is not defined`

Выполни:
```bash
cd back && go run main.go
```

Это пересоздаст таблицы с правильной схемой.

### ❌ Изображения не грузятся

1. Убедись что пути в БД правильные:
   ```sql
   SELECT name, image_url FROM menu_items LIMIT 5;
   ```

2. Если используешь относительные пути, они должны быть относительно HTML
3. Или используй полные URL: `https://example.com/images/Оливье.jpg`

### ❌ Корзина не работает

1. Проверь что все JS файлы подключены в `menu.html`
2. Открой консоль (F12) и проверь на ошибки
3. Убедись что `menuLoader.js` вызывается первым (в `main.js`)

## 📦 Файлы проекта

```
mdk2-Bd-menu/
├── back/
│   ├── main.go
│   ├── models/models.go (исправлено)
│   ├── db/
│   │   ├── db.go (исправлено)
│   │   └── postgres.go (исправлено)
│   └── handle/handle.go (готов)
│
├── js/
│   ├── menuLoader.js (НОВЫЙ ⭐)
│   ├── main.js (обновлено)
│   ├── cart.js
│   ├── cartUI.js
│   ├── dishModal.js
│   ├── orderModal.js
│   └── ... (остальные файлы)
│
├── menu.html (обновлено)
├── insert_menu.sql (помощник)
├── SETUP.md (инструкции)
└── MENU_MIGRATION.md (этот файл)
```

## ✅ Чек-лист перед продакшеном

- [ ] URL API верный в `menuLoader.js`
- [ ] БД заполнена данными (`insert_menu.sql` выполнен)
- [ ] CORS настроена правильно в `handle.go`
- [ ] Образы загружаются корректно
- [ ] Корзина работает (добавление, удаление, расчет суммы)
- [ ] Модали работают (клик на товар, открытие деталей)
- [ ] Категории переключаются корректно
- [ ] Нет ошибок в консоли браузера

## 🔐 Безопасность

- Убедись что `DATABASE_URL` не опубликована в Git
- Проверь что CORS Origin ограничен твоим доменом
- Используй HTTPS в продакшене

---

**Готово!** Меню теперь полностью управляется с бэка. Добавляй/удаляй товары прямо в БД.
