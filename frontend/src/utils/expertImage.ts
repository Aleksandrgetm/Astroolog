// Resolve legacy portrait URLs from the catalog without changing API data.
export function expertImage(source: string): string {
  return ['/images/Hero-new.png', '/images/expert-placeholder.webp', '/images/hero/hero-main.webp', '/images/hero.png'].includes(source)
    ? '/images/optimized/expert-1672.webp'
    : source
}
