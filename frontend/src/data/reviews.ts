// Stable IDs connect full, localized client stories to their i18n content.
export const reviews = Array.from({ length: 8 }, (_, index) => ({
  id: `review-${index + 1}`,
  textKey: `reviews.items.review${index + 1}`,
}))
