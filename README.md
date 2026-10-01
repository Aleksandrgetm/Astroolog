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

PostgreSQL — источник услуг и вариантов записи. `services` содержит группу, локализованные описания и изображение; `booking_options` — формат, цену в центах EUR, активность и признак дополнения `is_addon`. Публичная `Service.price` вычисляется как MIN активных основных вариантов (без дополнений); сохранённое поле `services.price` — legacy. NULL означает неизвестную цену, 0 — бесплатно.

Согласованные группы: бесплатное знакомство, отдельные разборы (€50, опциональная онлайн-встреча +€20), полный разбор (€250 / с двумя встречами €350), сопровождение (€50/150/350/600). Бонус нельзя выбрать без основного разбора. Четыре направления главной сохранены. Архивирование через `is_active`; FK используют RESTRICT, не каскадное удаление.

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

POST: `name`, `email`, `phone` (необязательно), `message`, `consent:true`, `language` (`ru/lv/en`). Для записи нужны `service_id`, `service_type` (код основного варианта), `price` (EUR cents либо NULL); для выбранного бонуса — `bonus_code` и `bonus_price`. Бонус должен быть активным дополнением той же услуги. Клиентский total игнорируется. Новая форма не передаёт `preferred_date`; новые заявки сохраняют NULL, исторические даты остаются в БД. Для вопроса message обязателен. Максимальное тело 16 KB. Телефон: 7–15 цифр, допускается международное форматирование.

Обязателен заголовок `Idempotency-Key`: 16–128 букв/цифр/`_`/`-`. Frontend создаёт UUID для конкретного payload, сохраняет его при сетевом retry. Повтор одинакового запроса возвращает успех без второй строки; другой payload с тем же ключом — 409 `idempotency_conflict`. Claim и заявка сохраняются одной транзакцией. Записи ключей пока не удаляются автоматически; будущая политика retention должна учитывать максимальное окно retry и удаление персональных данных. Fingerprint хранит SHA-256, не исходный текст.

Ошибки имеют локализуемый `code`; `stale_price` — 409. Цена/активность перепроверяются сервером под блокировкой строки. Бронирование сохраняет snapshot названий группы, варианта, формата, языка и валюты. `price` — основная цена, `bonus_code`, `bonus_title`, `bonus_price` — дополнение, `total_price` — итог. Ответ POST содержит сохранённый `booking` для подготовки текста мессенджера. Старые snapshot и NULL не переписываются.

`contact_requests`: статусы `new`, `in_progress`, `completed`, `archived`; `updated_at`, закрытый `internal_note`. Для старых обращений статус new. Бронирования сохраняют прежние статусы. В форме записи пользователь выбирает WhatsApp/Telegram до сохранения. После успешного POST открывается мессенджер; сообщение пользователь отправляет самостоятельно. Ошибка ссылки не отменяет сохранение. Повтор открытия ссылки не делает POST; выбор другого мессенджера с тем же payload использует прежний idempotency key и заявку. В sessionStorage хранится только хеш payload и случайный ключ, без полей формы. Форма вопроса сохраняет прежний flow.

## Миграции

`database.Migrate` выполняет транзакционную версионированную миграцию с PostgreSQL advisory lock. Маркеры хранятся в `catalog_revisions`.

- `schema-pre-admin-v1`: preflight отрицательных цен/осиротевших вариантов; привязка только двух известных legacy-кодов к `personal-matrix`; FK/NOT NULL/CHECK; workflow обращений, snapshot бронирования, таблица `submission_keys`; деактивация точных demo-отзывов.
- `catalog-bootstrap-v2`: только для пустого каталога — согласованный JSON `catalog_20260913.json`. Не перезаписывает непустую БД.
- `booking-bonus-20260926-v1`: деактивация исторического `reading-call-40`, дополнение `online-meeting-40` за 2000 cents, новые snapshot-поля и связь с idempotency key. Описания обновляются только при точном совпадении с bootstrap-текстом. Повторный запуск сохраняет admin-managed цены и тексты.
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

## Направления работы

Четыре editorial направления отделены от коммерческого каталога. Home cards ведут на `/directions/potential`, `/directions/purpose-money`, `/directions/relationships`, `/directions/child-matrix`; для LV/EN применяются обычные prefixes.

Один `DirectionPage.vue` получает направления и переводы через CMS API. `src/data/directions.json` и исходные i18n-файлы сохранены как bootstrap/fallback; текущий контент редактируется в `/admin/directions`.

Все направления связаны с `personal-matrix`: potential → personality, purpose-money → finances, relationships → relationships, child-matrix → child-matrix. Каждое также предлагает опциональный бонус `online-meeting-40` (+€20). Цены и доступность берутся из Pinia/API, editorial config не содержит цен. Для темы предназначения есть дополнительная ссылка на полный разбор. CTA сохраняет язык и выбирает соответствующий option в существующей форме.

Prerender и sitemap читают активные направления CMS API: добавлены 12 URL. Сейчас генерируются 42 контентные HTML + 3 локализованные 404; sitemap содержит 39 URL (Privacy исключена). Для direction social metadata используется изображение соответствующего направления. Breadcrumb «Направления» ведёт на Home#directions.

Дополнительная проверка: `PLAYWRIGHT_CHANNEL=chrome npm run test:directions` в frontend (production-router на TEST_SITE_URL, default 8087). Проверяет 60 сочетаний direction/язык/ширина, Home links, booking selection, Back, актуальную цену из API и fallback при недоступности каталога. Формы тестом не отправляются.

## Admin Panel и CMS

`/admin/login` → `/admin`. Административные маршруты не имеют языкового префикса и используют отдельный адаптивный layout. Интерфейс на русском; редакторы контента содержат вкладки RU/LV/EN. Пустые переводы явно видны в редакторе; публичный сайт может использовать RU fallback. Изображения сертификатов: выбранный язык → EN → RU.

Первый администратор создаётся вручную в интерактивном терминале, без default account:

```sh
cd backend
go run ./cmd/admin create
```

CLI запрашивает email и пароль дважды; пароль не отображается, не передаётся в командной строке и хранится только как bcrypt hash. CLI использует ту же конфигурацию PostgreSQL, что и сервер. Нельзя запускать интеграционные тесты с рабочей БД: тестовый аккаунт создаётся исключительно при явном `TEST_DB_NAME` с суффиксом `_test`.

Доступны: контент существующих секций, шаги процесса, группы услуг и варианты/бонусы, направления и связи с каталогом, отзывы, сертификаты, медиатека, заявки/вопросы, контакты, SEO, журнал действий и публикации. Кнопки «Архивировать»/«Восстановить» меняют доступность, не удаляя исторические заявки. Суммы вводятся в EUR и сохраняются целыми minor units; цена группы вычисляется из активных основных вариантов. Оригинальные сообщения и коммерческие snapshots в заявках не редактируются. Внутренняя заметка и рабочий статус отделены от исторического статуса booking.

Контент хранится в отдельных Page/PageSection, Direction, ProcessStep, Review, Certificate и SiteSettings. JSONB в переводах секций ограничен whitelist полей конкретной существующей секции; HTML/CSS/JS не принимаются. Коды услуг, вариантов и URL существующих направлений стабильны. Дизайн и структура формы контролируются кодом. Ранее удалённые информационные блоки направлений не возвращены.

Миграция `admin-cms-20260926-v1` выполняется один раз под advisory lock. Она переносит существующие тексты/8 отзывов/сертификаты/направления, регистрирует исходные изображения и сохраняет предыдущие бронирования. Повторный старт не перезаписывает изменения администратора. Seed в `backend/internal/database/cms_seed.json` является начальным снимком, не текущим хранилищем контента.

### Сессии и защита

Серверные сессии действуют 12 часов; в БД хранится hash токена. Cookie HttpOnly, SameSite=Strict, Path=/; production использует Secure и префикс `__Host-`. Logout инвалидирует сессию. CSRF token передаётся автоматически API-клиентом, Origin проверяется для изменяющих запросов. Login ограничивается по IP/email. Все административные API, кроме login, проверяют сессию, активность аккаунта и роль admin. Auth tokens не хранятся в localStorage. Админка/API исключены из индексации; robots не заменяет авторизацию.

### Production-конфигурация

Помимо существующих DB_*:

- `APP_ENV=production`, `APP_ORIGIN=https://ваш-домен` — точный HTTPS origin админки, без пути/query; публичный сайт и админка работают same-origin.
- `VITE_SITE_URL=https://ваш-домен` — canonical domain для сборки.
- `SITE_DIST_DIR=/absolute/path/frontend/dist` — начальная проверенная сборка.
- `CMS_FRONTEND_DIR=/absolute/path/frontend` — директория с package.json и установленными production build dependencies, включая devDependencies, необходимыми сборке.
- `CMS_RELEASE_ROOT=/persistent/path/releases` — writable directory для публикаций; одновременно доступна процессу Go.
- `PRERENDER_API_URL=http://127.0.0.1:8080/api` — внутренний API для сборки.
- `MEDIA_ROOT=/persistent/path/media` — постоянное хранилище загрузок, доступное на запись Go. Не размещайте его внутри временной директории релиза.
- Установленные Node/npm и Playwright Chromium; `PLAYWRIGHT_CHANNEL=chrome` допустим при наличии Chrome. Первый build выполняется после запуска API и миграции.

HTTPS должен завершаться на доверенном reverse proxy. Не открывайте PostgreSQL и служебный HTTP порт в интернет. Сервер по умолчанию не доверяет X-Forwarded-For; если нужен per-client IP rate limit за proxy, настройте доверенные proxy адреса явно в инфраструктуре/Go, не доверяйте произвольным заголовкам.

«Сохранить» обновляет CMS/API. «Опубликовать» собирает HTML, SEO и sitemap из актуального API в отдельную директорию и атомарно переключает symlink `current`. Пока идёт сборка, текущий релиз остаётся доступен. При ошибке прежняя версия сохраняется; статус отображается в Обзоре. Публикация не выполняет shell-команды из введённого контента. Уже открытые страницы могут получать content-hashed assets из сохранённых предыдущих релизов. Не удаляйте их сразу после публикации.

Изменения API-контента видны клиентскому Vue после перезагрузки уже до публикации; это не отдельный draft/approval workflow. Чтобы статический HTML и поисковики получили новые значения, публикуйте после завершения правок.

### Изображения и резервные копии

Загрузка принимает JPEG/PNG/WebP/AVIF по содержимому, до 10 MB и 20 MP. SVG/HTML/исполняемые файлы не принимаются. Изображения re-encode в WebP, EXIF удаляется, JPEG orientation применяется; варианты 640/960/1280/full создаются без увеличения исходника. Storage interface допускает будущую замену local filesystem на S3/R2. Используемое изображение нельзя удалить; исходные изображения проекта защищены от удаления через CMS.

Резервируйте PostgreSQL и MEDIA_ROOT вместе, отдельно сохраняйте исходный frontend/build configuration. Перед миграцией делайте `pg_dump`; восстановление сначала проверяйте на отдельной БД. Сессии и login-attempt записи можно удалять регламентным заданием по сроку давности, audit и заявки — только согласно согласованной политике хранения. Production домен, TLS, backups и фактический deploy на внешнем сервере настраиваются владельцем инфраструктуры.

### Проверки CMS

```sh
cd backend
TEST_DB_NAME=astroolog_admin_test go test ./...
go vet ./...
# В frontend, против отдельного локального сервера с тестовой БД:
TEST_SITE_URL=http://127.0.0.1:8086 npm run test:admin
# Только на тестовом сервере с CMS_RELEASE_ROOT в отдельном каталоге:
TEST_SITE_URL=http://127.0.0.1:8086 npm run test:admin:publish
```

Admin browser regression проверяет редакторы и списки на 390/430/768/1440 px, вход/выход, отсутствие публичного layout/overflow, локальный предпросмотр без запросов записи, предупреждение при уходе, сохранение и восстановление контента, языковые вкладки и медиатеку по запросу. DB tests проверяют RBAC, expiration/logout, CSRF, throttling, media upload/reference deletion и неизменность исходного сообщения/цены заявки.

### Редактирование контента

В «Контент сайта» выберите страницу и карточку секции. Редакторы открываются отдельной страницей, например `/admin/content/home/hero` или `/admin/content/home/how-it-works`. Поля объединены по смыслу, рядом доступен локальный предпросмотр. На экранах от 1200 px он закреплён справа; на меньших — под формой. Изменения не отправляются в API до явного сохранения. Закреплённая панель сообщает о несохранённых изменениях, уход со страницы требует подтверждения.

RU/LV/EN редактируются раздельно с индикатором заполненности; имя специалиста отображается вместо шаблонной переменной. Общие тексты и контактные данные отмечены пояснениями. Карточки услуг, направлений, отзывов и сертификатов ведут к отдельным редакторам; варианты услуги открываются из её карточки. Контакты находятся в настройках. SEO-редактор показывает пример поискового сниппета, счётчики символов и автоматически сформированные адреса.

Медиатека загружается только при открытии, по 12 файлов. Доступны поиск, фильтры, предпросмотр и ссылки на места использования. Замена изображения меняет выбор в текущем локальном черновике; её нужно сохранить. Используемые файлы защищены от удаления.

Для редакторов добавлено чтение одной записи `GET /api/admin/entities/:kind/:id`, а медиатека получила поиск, фильтр и сведения о местах использования. Существующие схемы БД и публикация сохранены. Исправлена проверка сохранения SEO: неизменяемые адреса страниц, включая `/`, больше не проверяются как коды услуг.

### Сброс пароля администратора и Origin

Из каталога `backend` выполните `go run ./cmd/admin reset-password`. Команда запрашивает email существующего администратора и дважды новый пароль в локальном терминале без отображения символов. Требования совпадают с `create`: 12–72 байта, bcrypt cost 12. Новый пользователь не создаётся, email/роль/активность не меняются. Изменение хеша, удаление всех сессий пользователя и аудит `password_reset` выполняются атомарно. Аудит содержит только источник `local_cli` и идентификатор затронутого аккаунта, без пароля/хеша. Сам запуск команды без завершённого ввода пароль не меняет.

Admin API читает `Config.AdminAllowedOrigins` из `ADMIN_ALLOWED_ORIGINS` — списка точных origins через запятую с удалением окружающих пробелов. Для локального frontend добавьте в свой непрокоммиченный `backend/.env`:

```dotenv
ADMIN_ALLOWED_ORIGINS=http://localhost:5174,http://127.0.0.1:5174
```

Явный список имеет приоритет над `APP_ORIGIN`. Если список отсутствует, используется настроенный `APP_ORIGIN`; если оба значения отсутствуют, только в development разрешены localhost/127.0.0.1:5174. Production по-прежнему требует точный HTTPS `APP_ORIGIN` и не получает localhost fallback. Origin сравнивается целиком, без wildcard и доверия заголовку Host. CSRF и свойства session cookies сохранены. Для изменения списка перезапустите backend.

Проверка CLI с записью в БД: `TEST_DB_NAME=astroolog_admin_final_test go test ./cmd/admin` — только с отдельной тестовой базой. Реальные пароли тесты не меняют.

### Legal documents and cookie consent

Privacy, Terms and Cookies are available in RU/LV/EN at `/privacy`, `/terms`, `/cookies` (with `/lv` and `/en` prefixes). Edit their named sections under **Контент сайта** in the existing CMS. Text remains escaped; blank lines separate paragraphs and lines starting with `- ` form lists. The additional 40-minute meeting price comes from the existing `online-meeting-40` catalog option.

Bundled fallback content lives in `frontend/src/legal/documents.json`. After changing that source, run `node frontend/scripts/sync-legal-seed.mjs`; `--check` verifies parity with the backend seed. The versioned migration adds the new sections without overwriting previous CMS content or later edits. Restart the backend to apply it to an existing installation, then use the existing publication workflow for production HTML.

Consent is managed by `frontend/src/legal/useCookieConsent.ts` and `consent.ts`. The optional integration registry is empty. New integrations must declare a category and start/stop handlers; their loader must only run through the consent controller. Bump `COOKIE_CONSENT_VERSION` when categories/policy materially change. Update the factual inventory and translations alongside technical changes. Technical cookie sections are read-only in CMS and rendered from implementation-maintained content. Cookie selection is separate from mandatory form privacy consent.

Client confirmation still needed before production: enquiry retention period; business registration details if applicable; production domain; hosting and email providers; intended payment provider; international-transfer arrangements. The implementation does not assert the DOCX's example 12-month retention period. Google Fonts is currently requested externally; WhatsApp/Telegram are opened by user action. The booking form records an enquiry and does not collect payment or record acceptance of Terms/early performance before expiry of a withdrawal period; agree that paid-order flow with the client before adding a separate acceptance step.

Validation: `npm run test:consent`, `npm run test:legal` (isolated local server), existing frontend tests, production build/prerender, and backend tests/vet. Browser integration tests must use an isolated database, never the client's working data.
