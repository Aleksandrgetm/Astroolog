import entries from './directions.json'
// Editorial topics are separate from Service/BookingOption; prices always come from the API.
export interface Direction {
  id: string
  slug: string
  contentKey: string
  image: string
  width: number
  height: number
  imagePosition: string
  relatedServiceSlug: string
  relatedBookingOptionCodes: string[]
}
export const directions: readonly Direction[] = entries
export const findDirection = (slug: string) => directions.find(direction => direction.slug === slug)
