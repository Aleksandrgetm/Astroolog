import { content } from '../stores/content'
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
export const findDirection = (slug: string): Direction | undefined => {
 if(!content.ready)return directions.find(direction=>direction.slug===slug)
 const d=content.directions.find(d=>d.slug===slug);if(!d)return undefined
 const original=directions.find(item=>item.slug===slug)
 return {id:d.slug,slug:d.slug,contentKey:d.content_key||d.slug,image:d.image,width:original?.width||1536,height:original?.height||1024,imagePosition:original?.imagePosition||'center',relatedServiceSlug:d.service_slug||'',relatedBookingOptionCodes:d.options.map((o:{code:string})=>o.code)}
}
