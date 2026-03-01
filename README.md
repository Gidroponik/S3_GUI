<p align="center">
  <img src="s3_gui.png" width="200" alt="S3 GUI Logo"/>
</p>

<h1 align="center">S3 Bucket Manager</h1>

<p align="center">
  Desktop-приложение для управления S3-совместимыми хранилищами.<br/>
  Красивый, быстрый, без лишнего.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white" alt="Go"/>
  <img src="https://img.shields.io/badge/Wails-v2-DB3552?logo=wails&logoColor=white" alt="Wails"/>
  <img src="https://img.shields.io/badge/Svelte-4-FF3E00?logo=svelte&logoColor=white" alt="Svelte"/>
  <img src="https://img.shields.io/badge/Tailwind-4-06B6D4?logo=tailwindcss&logoColor=white" alt="Tailwind"/>
</p>

<p align="center">
  <a href="https://github.com/Gidroponik/S3_GUI/releases/tag/v1.0.0">
    <img src="https://img.shields.io/badge/Download-v1.0.0-28A745?style=for-the-badge&logo=github&logoColor=white" alt="Download"/>
  </a>
</p>

<p align="center">
  <a href="https://github.com/Gidroponik/S3_GUI/releases/tag/v1.0.0">
    <img src="https://img.shields.io/badge/Windows-amd64-0078D6?style=flat-square&logo=windows&logoColor=white" alt="Windows"/>
  </a>
  <a href="https://github.com/Gidroponik/S3_GUI/releases/tag/v1.0.0">
    <img src="https://img.shields.io/badge/macOS-Universal-000000?style=flat-square&logo=apple&logoColor=white" alt="macOS"/>
  </a>
</p>

---

## Скачать

Готовые сборки — не нужно ничего компилировать:

| Платформа | Файл | Архитектура |
|-----------|------|-------------|
| <img src="https://img.shields.io/badge/-Windows-0078D6?logo=windows&logoColor=white&style=flat-square" /> | [S3BucketGUI.exe](https://github.com/Gidroponik/S3_GUI/releases/tag/v1.0.0) | amd64 |
| <img src="https://img.shields.io/badge/-macOS-000000?logo=apple&logoColor=white&style=flat-square" /> | [S3BucketGUI.dmg](https://github.com/Gidroponik/S3_GUI/releases/tag/v1.0.0) | Universal (Intel + Apple Silicon) |

> **[Все релизы →](https://github.com/Gidroponik/S3_GUI/releases)**

---

## Предыстория

Искал готовое desktop-решение для работы с S3-хранилищами. Перепробовал несколько вариантов — одни перегружены лишним функционалом, другие выглядят как из 2005 года, третьи стоят денег за базовые вещи. В итоге ни одно не подошло, и я решил сделать своё: лёгкое, с нормальным UI, ровно с тем набором функций, который реально нужен.

## Возможности

### Подключения
- Поддержка любых S3-совместимых хранилищ (MinIO, Ceph, Backblaze B2, AWS S3 и др.)
- Несколько профилей подключений с быстрым переключением
- SSL и Path-Style конфигурация
- Тест соединения перед сохранением
- Экспорт / импорт профилей (`.s3b`) для переноса между устройствами
- Учётные данные шифруются AES-256-GCM с ключом, привязанным к машине

### Файловый менеджер
- Навигация по папкам с breadcrumb-путём
- Загрузка файлов и папок (диалог + drag & drop)
- Скачивание файлов и целых директорий
- Создание папок
- Удаление файлов и папок (рекурсивно, с подтверждением)
- Мультивыбор с Ctrl / Shift
- Контекстное меню (ПКМ) с основными действиями
- Копирование Presigned URL (с TTL) и Direct URL
- Сортировка по имени, размеру, дате

### Трансферы
- Параллельные загрузки / скачивания (настраивается от 1 до 10)
- Прогресс-бар в реальном времени для каждого файла
- Отмена отдельных трансферов
- Очистка завершённых

### UI
- Тёмная тема
- Нативное окно (без Electron, ~17 MB)
- Горячие клавиши: `Ctrl+R` обновить, `Ctrl+A` выбрать все, `Delete` удалить, `Backspace` назад
- Toast-уведомления
- Модальные диалоги для настроек и подключений

## Стек

| Слой | Технология |
|------|-----------|
| Backend | Go + AWS SDK v2 |
| Frontend | Svelte 4 + Tailwind CSS 4 |
| Desktop | Wails v2 (WebView2) |
| Шифрование | AES-256-GCM + PBKDF2 (machine-bound key) |

## Сборка из исходников

### Требования
- [Go](https://go.dev/) 1.21+
- [Node.js](https://nodejs.org/) 18+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2

### Шаги

```bash
# Установка Wails CLI (если ещё не установлен)
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Клонирование
git clone https://github.com/Gidroponik/S3_GUI.git
cd S3_GUI

# Сборка
wails build
```

Готовый бинарник: `build/bin/S3BucketGUI.exe`

### Dev-режим

```bash
wails dev
```

## Хранение данных

Все данные хранятся локально:

```
%APPDATA%/S3BucketGUI/          # Windows
~/Library/Application Support/S3BucketGUI/  # macOS
├── connections.dat   # профили (AES-256-GCM)
└── settings.json     # настройки (plain JSON)
```

## Лицензия

MIT
