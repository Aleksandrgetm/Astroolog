# PRE-ADMIN + Production Technical SEO — отчёт

Дата: 16 сентября 2026. Стек сохранён, redesign/Admin UI не выполнялись. Commit, push и deploy не выполнялись. Рабочая БД не изменялась: миграции проверены на её изолированной копии и на пустой тестовой БД. При запуске новой версии backend миграция применяется автоматически.

## A. Исправленные проблемы аудита

Закрыты tracked env, публичные demo testimonials, телефон из одних знаков, UTC date boundary, небезопасное формирование Telegram URL после POST, вечный cache каталога, независимая публичная цена группы, варианты неактивных групп, повторные POST, отсутствие snapshots/workflow, повторные content migrations, потери query, якоря и Back/Forward. Убраны технически неверные публичные утверждения Privacy и подтверждённый starter code.

## B. Security/env

`backend/.env` и `frontend/.env` сняты с Git tracking; локальные файлы сохранены. Оба `.env.example` безопасны, реальные env игнорируются. `backend/.env` найден в истории (b991e7a, 2026-09-09). Необходимо сменить использовавшийся пароль PostgreSQL в БД и окружениях; значения не выводились, история не переписывалась. Production не читает dotenv, требует DB environment и проверяемое TLS-соединение. Development PostgreSQL привязан к loopback. SQL-логирование параметризовано, HTTP timeout и DB pool ограничены. HTTP redirects нормализуют путь без переходов на внешний host; JSON-LD безопасно экранирует HTML-разделитель при prerender.

## C. Database migrations

Версионированная транзакция, advisory lock против параллельных startup, preflight отрицательных цен и orphan options. Два известных legacy option связаны с personal-matrix. FK RESTRICT и NOT NULL для service_id, CHECK неотрицательных цен. Добавлены workflow/snapshot/idempotency поля. Неизвестные старые цены остаются NULL.

На копии после миграции: 19 бронирований, 18 обращений, 16 NULL-цен — сохранены. Три demo-отзыва деактивированы, не удалены. Integration tests создают только собственные записи в тестовых БД и удаляют их после проверки.

## D. Router/scroll/query

Языковые prefixes и обратная совместимость query. Якорь ждёт DOM/асинхронный контент через MutationObserver, шрифты и предшествующие изображения; учитывается sticky header. Back возвращает savedPosition после готовности страницы. Устаревшая отложенная прокрутка не вмешивается в следующую навигацию. Contacts меняет только booking, сохраняя остальные параметры/hash. Смена языка не размонтирует текущую форму.

## E. Forms/validation/idempotency

Телефон 7–15 цифр с международным форматированием на обеих сторонах. Europe/Riga для минимальной даты. UUID для конкретного payload, повтор неизменённого запроса использует тот же ключ. Backend атомарно создаёт claim и заявку; конкурентные повторы возвращают успех без дублей. Конфликт payload — 409. Messenger helper больше не отменяет успешное сохранение. Сообщения никому автоматически не отправлялись.

## F. Catalog preparation for Admin

Публичная цена группы — MIN активных BookingOptions. Активность родителя обязательна. TTL 60 секунд на следующий load, force refresh/invalidate, без polling. Stale-price приводит к обновлению каталога и явному сообщению. Исторические snapshots защищают смысл уже созданных заявок. Admin UI/auth не создавались.

## G. Privacy

Убраны публичные «проект документа»/«сведения уточняются». Описание хранения языка исправлено на URL. TODO юридического имени, регистрации, controller contacts, оснований, получателей и retention оставлен в source/README; сведения не выдуманы. Privacy получает noindex, follow и исключена из sitemap.

## H. Removed dead/demo code

Удалены HelloWorld, starter Vue/Vite SVG, starter hero asset, неиспользуемый public icons.svg, VSelect registration/defaults, frontend testimonial API/types и публичный backend endpoint. Оригиналы бизнес-фото и сертификатов сохранены. MDI оставлен, поскольку реально используется.

## I. SEO architecture

Центральный helper для metadata, отдельные языковые пути, prerender из действующего API, manifest для production HTTP routing. Нет скрытых keyword paragraphs и искусственных рейтингов.

## J. Новые URL RU/LV/EN

RU `/`, `/about`, `/services`, `/services/:slug`, `/reviews`, `/contacts`, `/privacy`; LV и EN — те же пути с `/lv` и `/en`. Главные языковые страницы `/lv/`, `/en/`. Slugs стабильны, browser-language redirect отсутствует. Legacy `?lang` сохраняет остальные query.

## K. Canonical

Self canonical, один на страницу, абсолютный HTTPS origin из обязательного VITE_SITE_URL. Query/hash исключены. LV/EN не canonicalize на RU. Неизвестный production domain не подставлен.

## L. Hreflang

На локализованных страницах reciprocal ru/lv/en/x-default; x-default указывает на RU. Те же alternates в sitemap. HTML lang синхронизирован с URL.

## M. Sitemap

Генерируется в dist из реальных routes и активных API services. 27 indexable URL при четырёх текущих услугах. Без privacy, query, redirect, growth-point, API, 404 и fake lastmod. После публикации изменений каталога нужен rebuild.

## N. robots.txt

Генерируется с реальным site origin и Sitemap URL, разрешает публичные CSS/JS/images. Будущая админка должна защищаться auth/RBAC, а не robots.txt.

## O. Structured data

Person на Home/About: подтверждённое имя, URL, фото, локализованный jobTitle. BreadcrumbList на detail соответствует видимой навигации. Нет fake address, awards, sameAs, aggregateRating или review stars.

## P. Open Graph/Twitter

Title, description, type, absolute URL/image, размеры, alt, locale/alternate locales; Twitter summary_large_image/title/description/image. Verification meta выводится только при заданном реальном токене.

## Q. 404/noindex/redirect handling

Go `internal/web` обслуживает manifest известного dist. Неизвестные пути/услуги получают HTTP 404 и noindex; API не попадает в SPA fallback. Growth-point и legacy language query получают HTTP 301. Inactive detail дополнительно проверяется по БД. Production требует SITE_DIST_DIR и TLS reverse proxy; правила описаны в README.

## R. Prerender

Vite + Playwright, без Nuxt/смены архитектуры. 30 готовых HTML и три локализованных 404. Catalog API читается во время build; недоступность/невалидный каталог останавливает сборку. HTML содержит H1, текст, metadata и релевантный JSON-LD без JS. Затем Vue выполняет обычный client mount — это не SSR hydration. При публикации из будущей админки нужен workflow rebuild и атомарного dist deploy.

## S. Images/Core Web Vitals

Hero 1672px: примерно 1620 → 84 КБ; responsive 640/960/1280/1672. Сертификаты: 800/1200/full WebP, оригиналы сохранены. Hero eager/high priority, размеры заданы, ниже первого экрана lazy/async. Шрифты подключены stylesheet/preconnect вместо CSS @import. Семантика H2/H3 улучшена без изменения размеров кнопок/макета. Реальные field CWV пока не измеримы без production/CrUX; гарантии показателей не заявляются.

## T. Tests

- TypeScript + Vite + prerender: проходят с тестовым HTTPS origin `https://example.invalid`, без публикации.
- Build без VITE_SITE_URL: ожидаемый fail, защита проверена.
- test:i18n, test:whatsapp, test:pre-admin: проходят.
- test:seo: 30 HTML без JavaScript, H1, metadata, canonical, hreflang, lang, JSON-LD, alt, sitemap.
- Responsive: 150 сочетаний 30 страниц × 390/430/768/1024/1440, без horizontal overflow/page errors.
- Interaction: Home→anchor, direct anchor, Back/Forward, заполненная форма→LV, query/hash, повтор POST с тем же ключом после сетевой ошибки, mobile menu/focus/Escape.
- Go test ./..., go vet ./...: проходят. Integration на копии и пустой БД: проходят; шесть конкурентных повторов создают одну запись, snapshot/stale price/активность/вычисленная цена/restart сохраняют корректность.
- HTTP routing tests: 301/404 и закрытый manifest/demo API.
- Скриншоты Hero/сертификатов 390 и 1440 проверены визуально; изображения загружаются, лицо не обрезано.
- iPhone/Android hardware и Safari/Windows отдельно не запускались; responsive проверялся в Chrome.

## U. Созданные migrations

`schema-pre-admin-v1` в database/migrate.go и empty-catalog bootstrap `catalog-bootstrap-v2` в database/catalog.go. Маркеры хранятся в catalog_revisions. Следующие изменения схемы требуют нового version step. Автоматического destructive rollback нет.

## V. Изменённые файлы

Полный список приведён в конце отчёта, включая новые scripts/tests/assets. Изменения остаются в working tree; staged только снятие реальных .env с tracking.

Тестовые БД `astroolog_pre_admin_test` и `astroolog_bootstrap_test` остаются локально; production deploy не выполнялся.

## W. Вручную перед production

Указать реальный VITE_SITE_URL, DB environment/TLS и PRERENDER_API_URL. Сменить ранее tracked пароль. Подтвердить юридические данные/retention. Сделать backup, применить миграцию новой версией backend. Установить Chromium для CI, собрать с реальным API, атомарно разместить dist, настроить HTTPS/reverse proxy/SITE_DIST_DIR и перезапуск. Проверить реальный домен/301/404. Добавить Search Console property/verification, отправить sitemap и выполнить URL Inspection. Текущий dist с example.invalid НЕ публиковать.

## X. Перед Admin Panel

Auth/session/CSRF/RBAC, bootstrap администратора, аудит действий, CRUD/archive, модель реальных отзывов и assets, retention/backup, publish→rebuild pipeline. Idempotency key retention пока не автоматизирован. Прямое изменение каталога без rebuild не обновляет статический HTML/sitemap. Наличие User.Role само по себе не является защитой API.

### Полный список файлов

- `.gitignore`
- `README.md`
- `backend/.env`
- `backend/.env.example`
- `backend/cmd/server/main.go`
- `backend/internal/config/config.go`
- `backend/internal/database/catalog.go`
- `backend/internal/database/database.go`
- `backend/internal/database/migrate.go`
- `backend/internal/handlers/handlers.go`
- `backend/internal/models/models.go`
- `backend/internal/repositories/repository.go`
- `backend/internal/routes/routes.go`
- `backend/internal/services/integration_test.go`
- `backend/internal/services/requests.go`
- `backend/internal/services/requests_test.go`
- `backend/internal/web/static.go`
- `backend/internal/web/static_test.go`
- `docker-compose.yml`
- `docs/pre-admin-production-report.md`
- `frontend/.env`
- `frontend/.env.example`
- `frontend/index.html`
- `frontend/package-lock.json`
- `frontend/package.json`
- `frontend/public/icons.svg`
- `frontend/public/images/optimized/expert-1280.webp`
- `frontend/public/images/optimized/expert-1672.webp`
- `frontend/public/images/optimized/expert-640.webp`
- `frontend/public/images/optimized/expert-960.webp`
- `frontend/public/images/optimized/icta-en-1200.webp`
- `frontend/public/images/optimized/icta-en-800.webp`
- `frontend/public/images/optimized/icta-en-full.webp`
- `frontend/public/images/optimized/icta-ru-1200.webp`
- `frontend/public/images/optimized/icta-ru-800.webp`
- `frontend/public/images/optimized/icta-ru-full.webp`
- `frontend/public/images/optimized/numerology-en-1200.webp`
- `frontend/public/images/optimized/numerology-en-800.webp`
- `frontend/public/images/optimized/numerology-en-full.webp`
- `frontend/public/images/optimized/numerology-lv-1200.webp`
- `frontend/public/images/optimized/numerology-lv-800.webp`
- `frontend/public/images/optimized/numerology-lv-full.webp`
- `frontend/public/images/optimized/numerology-ru-1200.webp`
- `frontend/public/images/optimized/numerology-ru-800.webp`
- `frontend/public/images/optimized/numerology-ru-full.webp`
- `frontend/scripts/prerender.mjs`
- `frontend/scripts/test-browser.mjs`
- `frontend/scripts/test-interactions.mjs`
- `frontend/scripts/test-seo.mjs`
- `frontend/src/App.vue`
- `frontend/src/assets/hero.png`
- `frontend/src/assets/vite.svg`
- `frontend/src/assets/vue.svg`
- `frontend/src/components/BookingServiceChoice.vue`
- `frontend/src/components/CertificateGallery.vue`
- `frontend/src/components/HelloWorld.vue`
- `frontend/src/components/HeroSection.vue`
- `frontend/src/components/LocaleLink.vue`
- `frontend/src/components/MobileNavigation.vue`
- `frontend/src/components/RequestForm.vue`
- `frontend/src/components/ServiceCards.vue`
- `frontend/src/components/ServiceOptions.vue`
- `frontend/src/data/certificates.ts`
- `frontend/src/i18n/index.ts`
- `frontend/src/i18n/locales/en.ts`
- `frontend/src/i18n/locales/lv.ts`
- `frontend/src/i18n/locales/ru.ts`
- `frontend/src/i18n/routing.ts`
- `frontend/src/i18n/seo.ts`
- `frontend/src/i18n/urls.ts`
- `frontend/src/layouts/SiteLayout.vue`
- `frontend/src/main.ts`
- `frontend/src/pages/AboutPage.vue`
- `frontend/src/pages/ContactsPage.vue`
- `frontend/src/pages/HomePage.vue`
- `frontend/src/pages/PrivacyPage.vue`
- `frontend/src/pages/ServicePage.vue`
- `frontend/src/plugins/vuetify.ts`
- `frontend/src/router/index.ts`
- `frontend/src/router/scroll.ts`
- `frontend/src/services/api.ts`
- `frontend/src/stores/catalog.ts`
- `frontend/src/style.css`
- `frontend/src/types/index.ts`
- `frontend/src/utils/expertImage.ts`
- `frontend/src/utils/telegram.ts`
- `frontend/src/utils/validation.ts`
- `frontend/tests/pre-admin.test.ts`
- `frontend/vite.config.ts`
