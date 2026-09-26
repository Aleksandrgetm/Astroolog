# Astroolog

Персональный сайт Елены Захаровой: Vue 3 / TypeScript / Vite / Vuetify / Vue Router / Pinia / vue-i18n; API на Go / Gin / GORM; PostgreSQL 17. Admin UI и авторизация администратора пока не реализованы.

## Development

1. `docker compose up -d` из корня. Compose предназначен только для разработки; PostgreSQL опубликован **на loopback** `127.0.0.1:5433`.
2. Создайте локальный `backend/.env` по `.env.example`, указав пароль своей development БД. Для существующего Docker volume смена `POSTGRES_PASSWORD` не меняет пароль внутри БД.
3. В `backend`: `go run ./cmd/server` (API на 8080).
4. В `frontend`: `npm ci`, затем `npm run dev -- --host 127.0.0.1`. Vite проксирует `/api` на backend.

Реальные `.env` игнорируются Git, примеры отслеживаются. `.env` был найден в истории: ранее использованный пароль БД необходимо ротировать в БД и окружении. История автоматически не переписывается.

## URL и локализация

RU — `/`, `/about`, `/services`, `/services/:slug`, `/reviews`, `/contacts`, `/privacy`.
LV — `/lv/`, `/lv/about`, `/lv/services` и соответствующие остальные пути.
EN — `/en/`, `/en/about`, `/en/services` и соответствующие остальные пути.

Slugs одинаковы во всех языках: `first-step`, `personal-matrix`, `full-matrix`, `path-to-self`. Legacy `?lang=lv/en` переносится в prefix, остальные query сохраняются. `/services/growth-point` перенаправляется на `/services` с сохранением языка. В production это HTTP 301, в Vite — совместимость через Router. Язык определяется URL, без localStorage и без автоматического перенаправления по языку браузера. Используйте `LocaleLink` для внутренних ссылок. Смена языка сохраняет query/hash и экземпляр текущей страницы/формы. Полное обновление страницы не сохраняет введённые данные.

## Каталог и отзывы

PostgreSQL — источник услуг и вариантов записи. `services` содержит группу, локализованные описания и изображение; `booking_options` — формат, цену в центах EUR и активность. Публичная `Service.price` вычисляется как MIN активных вариантов; сохранённое поле `services.price` — legacy, будущая админка не должна редактировать его как независимую публичную цену. NULL означает неизвестную цену, 0 — бесплатно.

Согласованные группы: бесплатное знакомство, отдельные разборы (€50 / с онлайн-встречей €70), полный разбор (€250 / с двумя встречами €350), сопровождение (€50/150/350/600). Четыре направления главной сохранены и ведут на соответствующие варианты. Архивирование через `is_active`; FK используют RESTRICT, не каскадное удаление.

Pinia: `load({force:true})`, `invalidate()`, TTL 60 секунд при следующем обращении к `load`; фонового polling нет. При `stale_price` форма обновляет каталог и просит проверить выбор. Сохранённая цена заявки не меняется.

Реальные отзывы пока находятся в `frontend/src/data/reviews.ts`, `src/i18n/reviews/*.json` и коротких цитатах главной. Публичный `/api/testimonials` удалён, три точных demo-записи деактивированы, исторические строки сохранены. Нет выдуманных рейтингов или Review schema. Будущую Review model следует проектировать с учётом этих текстов, языков и согласий на публикацию.

## API и формы

| Метод | Путь | Назначение |
|---|---|---|
| GET | `/api/health` | API/DB health |
| GET | `/api/services` | Активные группы с вычисленной минимальной ценой |
| GET | `/api/services/:slug` | Активная услуга или 404 |
| GET | `/api/booking-options` | Только активные варианты активных групп |
| POST | `/api/contact` | Обращение |
| POST | `/api/bookings` | Запись |

POST: `name`, `email`, `phone` (необязательно), `message`, `consent:true`, `language` (`ru/lv/en`). Для записи нужны `service_id`, `service_type` (код варианта), `price` (EUR cents либо NULL), необязательная `preferred_date` YYYY-MM-DD. Для вопроса message обязателен. Максимальное тело 16 KB. Телефон: 7–15 цифр, допускается международное форматирование. Дата сравнивается с текущим днём Europe/Riga, включая переходы летнего времени.

Обязателен заголовок `Idempotency-Key`: 16–128 букв/цифр/`_`/`-`. Frontend создаёт UUID для конкретного payload, сохраняет его при сетевом retry. Повтор одинакового запроса возвращает успех без второй строки; другой payload с тем же ключом — 409 `idempotency_conflict`. Claim и заявка сохраняются одной транзакцией. Записи ключей пока не удаляются автоматически; будущая политика retention должна учитывать максимальное окно retry и удаление персональных данных. Fingerprint хранит SHA-256, не исходный текст.

Ошибки имеют локализуемый `code`; `stale_price` — 409. Цена/активность перепроверяются сервером под блокировкой строки. Бронирование сохраняет snapshot названий группы, варианта, формата, языка, цены и валюты. Старые NULL не восстанавливаются выдуманными значениями.

`contact_requests`: статусы `new`, `in_progress`, `completed`, `archived`; `updated_at`, закрытый `internal_note`. Для старых обращений статус new. Бронирования сохраняют прежние статусы. Автоматической отправки email/сообщений нет: после сохранения пользователь может открыть мессенджер. Неисправная ссылка мессенджера не отменяет успешную заявку. Переход не отправляет сообщение автоматически.

## Миграции

`database.Migrate` выполняет транзакционную версионированную миграцию с PostgreSQL advisory lock. Маркеры хранятся в `catalog_revisions`.

- `schema-pre-admin-v1`: preflight отрицательных цен/осиротевших вариантов; привязка только двух известных legacy-кодов к `personal-matrix`; FK/NOT NULL/CHECK; workflow обращений, snapshot бронирования, таблица `submission_keys`; деактивация точных demo-отзывов.
- `catalog-bootstrap-v2`: только для пустого каталога — согласованный JSON `catalog_20260913.json`. Не перезаписывает непустую БД.
- Исторические файлы переводов/данных сохранены, но старые перезаписывающие content migrations больше не вызываются при startup. `SEED_DEMO` больше не используется.

Повторный запуск с применённым маркером не запускает AutoMigrate/seed и не меняет редакторский контент. Для следующей schema migration нужен новый явный version step; нельзя просто менять модели и ожидать автоматического изменения существующей БД. Перед применением: backup и проверка на копии. Миграция отказывается продолжать при неизвестных orphan/negative-price данных, транзакция откатывается. Автоматический down не предусмотрен; восстановление из проверенного backup.

## Production и prerender

Production DB получает `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE` из environment при `APP_ENV=production`. `.env` в этом режиме не загружается, все значения обязательны; `DB_SSLMODE=verify-full`, доверенный CA настроен на хосте. Не публикуйте порт production PostgreSQL; используйте приватную сеть/firewall. `APP_PORT` — порт API, `GIN_MODE=release`.

Frontend production build обязательно требует `VITE_SITE_URL` — реальный HTTPS origin без пути, query и credentials. Неизвестный домен не подставляется. `VITE_GOOGLE_SITE_VERIFICATION` необязателен, пустая meta не создаётся. `VITE_API_URL` обычно пустой (same-origin `/api`), при отдельном API нужен собственный безопасный CORS на gateway. `VITE_*` — публичные build-time значения, никогда не помещайте туда secrets.

В `frontend`:

```sh
npm ci
npx playwright install --with-deps chromium
# Установите VITE_SITE_URL и PRERENDER_API_URL через окружение CI.
npm run build
```

`PRERENDER_API_URL` — доступный сборке текущий API с `/api` в конце; development default `http://127.0.0.1:8080/api`. Если данные недоступны, сборка падает, не подменяя их старым seed. Можно использовать установленный Chrome через `PLAYWRIGHT_CHANNEL=chrome`.

Сборка сначала Vite, затем `scripts/prerender.mjs`: снимок действующего каталога, 42 HTML страницы (при четырёх активных услугах и четырёх направлениях) и три 404. HTML содержит контент и метаданные без JS. После загрузки JS Vue монтируется поверх снимка и продолжает обычную клиентскую навигацию; это **prerender с client mount, не Vue SSR hydration**. Данные форм в снимках отсутствуют. Для изменения каталога/отзывов нужен новый build и атомарная публикация dist. Будущая админка должна запускать rebuild после публикации; иначе статический HTML и sitemap могут отставать от API. Сервер дополнительно проверяет активность detail перед выдачей.

`npm run build:client` — development-only сборка без prerender; не является production SEO-артефактом. `https://example.invalid` допустим только как временная тестовая fixture, не как домен публикации.

Рекомендуемая минимальная выдача: Go backend с `SITE_DIST_DIR` (абсолютный путь к готовому `frontend/dist`) за TLS reverse proxy. `internal/web` читает `routes.json`, отдаёт только известные HTML routes, реальные 404 и legacy 301. API и assets обслуживаются отдельно, произвольный SPA fallback 200 для неизвестных URL отсутствует. На CDN нужно повторить те же правила; обычного `try_files ... /index.html` недостаточно. Proxy должен сохранять path/query и отправлять весь трафик к Go; настройте HTTPS, сжатие, лимиты запросов и безопасное логирование. Dist публикуется атомарно с перезапуском процесса для перечитывания manifest.

## SEO

Единый `src/i18n/seo.ts`: self canonical без query/hash, reciprocal ru/lv/en/x-default (RU), html lang, локализованные title/description, OG + Twitter, абсолютные URL изображений. Person только Home/About, BreadcrumbList только у видимых breadcrumbs услуг. Нет fake address/awards/ratings. Privacy и неизвестные страницы — `noindex, follow`.

`dist/sitemap.xml` генерируется из активных API services и реальных routes, содержит языковые alternates; исключает privacy, API, 404, growth-point, query duplicates. Fake lastmod нет. `dist/robots.txt` разрешает публичные ресурсы и указывает sitemap. Robots.txt не защищает будущую админку: ей нужны auth, права и noindex.

Hero и сертификаты используют WebP; оригиналы сохранены. Hero имеет srcset, размеры и high fetch priority; ниже первого экрана lazy/async. Google Fonts подключены stylesheet + preconnect, без CSS @import. MDI сохранён: реально используется Vuetify и интерфейсом. Большой CSS не переписывался ради риска регрессии.

### SEO / Google Search Console setup

После реального deploy: добавить property, подтвердить домен (предпочтительно DNS; либо предоставленный verification meta), отправить `/sitemap.xml`, проверить языковые URL через URL Inspection, запросить индексацию основных страниц. Проверить HTTP redirects/404 и canonical на реальном HTTPS домене. Property автоматически не создаётся. Изменение build-time env требует новой сборки.

## Privacy и ручные данные

Публичный текст описывает действующие формы и URL-язык без developer-заглушки. В source оставлен TODO: юридическое имя/registration/controller contacts, сроки хранения, основания и получатели обработки должны быть подтверждены владельцем перед production. Неизвестные сведения не выдуманы. Текущие подтверждённые контакты: +371 29 580 232, jelenabobrovska@gmail.com; имя/бренд не переводятся.

## Проверки

```sh
# frontend
npm run test:i18n
npm run test:whatsapp
npm run test:pre-admin
npm run build
npm run test:seo
PLAYWRIGHT_CHANNEL=chrome node scripts/test-browser.mjs
# backend
 go test ./...
 go vet ./...
# Только изолированная БД с суффиксом _test:
TEST_DB_NAME=astroolog_pre_admin_test go test ./...
```

Browser regression использует `TEST_SITE_URL` (default `http://127.0.0.1:8087`) с production-router и проверяет 390/430/768/1024/1440 для всех языков. SEO tests читают dist без выполнения JS. Интеграционные DB tests пропускаются без TEST_DB_NAME. Не направляйте тесты на рабочую БД.

## Перед Admin Panel

Нужны отдельный дизайн auth/session/CSRF и RBAC, безопасный bootstrap администратора, аудит действий, CRUD с архивированием и транзакционными проверками, политика retention/backup, модель реальных отзывов и управление assets, workflow publish/rebuild SEO. Поле role/password_hash в User само по себе не является реализованной защитой. Admin UI, логин, email-рассылки и deploy в этот этап не входят.


## Направления работы

Четыре editorial направления отделены от коммерческого каталога. Home cards ведут на `/directions/potential`, `/directions/purpose-money`, `/directions/relationships`, `/directions/child-matrix`; для LV/EN применяются обычные prefixes.

Один `DirectionPage.vue` использует конфигурацию `src/data/directions.json` через типизированный `directions.ts`. Поля: id, slug, contentKey, image/dimensions/position, relatedServiceSlug и relatedBookingOptionCodes. Контент хранится в `src/i18n/directions/{ru,lv,en}.json`.

Все направления связаны с `personal-matrix`: potential → personality, purpose-money → finances, relationships → relationships, child-matrix → child-matrix. Каждое также предлагает `reading-call-40`. Цены и доступность берутся из Pinia/API, editorial config не содержит цен. Для темы предназначения есть дополнительная ссылка на полный разбор. CTA сохраняет язык и выбирает соответствующий option в существующей форме.

Prerender и sitemap читают тот же directions.json: добавлены 12 URL. Сейчас генерируются 42 контентные HTML + 3 локализованные 404; sitemap содержит 39 URL (Privacy исключена). Для direction social metadata используется изображение соответствующего направления. Breadcrumb «Направления» ведёт на Home#directions.

Дополнительная проверка: `PLAYWRIGHT_CHANNEL=chrome npm run test:directions` в frontend (production-router на TEST_SITE_URL, default 8087). Проверяет 60 сочетаний direction/язык/ширина, Home links, booking selection, Back, актуальную цену из API и fallback при недоступности каталога. Формы тестом не отправляются.
