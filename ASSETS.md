# Временные изображения

Созданы встроенным ImageGen, затем преобразованы в WebP для сайта. Это демонстрационные изображения, а не фотографии реального клиента. Основные пути:

- `frontend/public/images/expert-placeholder.webp` — портрет, 116 KB.
- `frontend/public/images/journal-placeholder.webp` — блокнот и кофе, 268 KB.
- `frontend/public/images/relationships-placeholder.webp` — руки, отношения.
- `frontend/public/images/child-placeholder.webp` — ребёнок за рисованием.

Исходные PNG портрета и блокнота сохранены в `frontend/src/assets/source/` и не включаются в публичную сборку. Для замены достаточно загрузить согласованные фотографии по WebP-путям. Названия услуг и изображения каталога хранятся в таблице services.

## Промпты (встроенный ImageGen)

### Портрет
Use case: photorealistic-natural. Asset type: premium personal coaching website portrait placeholder. Create a warm editorial portrait photograph of a fictional female coach, early 40s, shoulder length wavy light brown hair, subtle glasses, beige linen blazer over ivory top. Seated at a light oak desk with an open cream notebook and ceramic coffee cup. Relaxed confident friendly expression, looking at camera, natural hands loosely resting on notebook. Warm ivory plaster interior, sunlit sheer curtains, subtle dried branch in vase to side. Waist-up portrait, centered subject, 4:5 vertical framing, generous space above head. Restrained cream sand caramel colors, natural skin texture, soft daylight, sophisticated magazine photography. No text, no collage, no logos. This is a placeholder fictional person, not a real client.

### Блокнот
Use case: photorealistic-natural. Asset type: premium coaching website still life photograph. A beautiful cream linen hardcover journal partly open, a slim antique brass pen diagonally across it, ceramic cup of coffee and dried delicate branches on softly rumpled ivory linen next to sunlit warm travertine table. Top down at gentle angle, refined editorial lifestyle photography, warm light beige caramel brown monochromatic colors, afternoon shadows, tactile real materials, calm spacious composition. Landscape 3:2. No text, no letters, no logos, no people.

### Отношения
Use case: photorealistic-natural. Asset type: relationship coaching website card. Close up editorial photograph of two adult people's hands gently holding each other over a light oak table, beige linen sleeves, soft warm window light, quiet intimacy, natural anatomy, elegant cream ivory sand palette matching a premium feminine coaching website. Beautiful shallow depth of field, tactile skin and linen textures. Landscape 3:2. No faces, no text, no collage, no logos.

### Ребёнок
Use case: photorealistic-natural. Asset type: family coaching website card. Warm editorial candid photograph of a young child about 7 years old with light brown hair wearing an ivory knit sweater drawing with a wooden pencil in a sketchbook at a warm wooden table. Side profile concentrating quietly, soft natural window daylight in a beige home, cream sand warm brown palette. Calm affectionate authentic mood. Landscape 3:2 medium close framing. No text, no logo, no collage.
