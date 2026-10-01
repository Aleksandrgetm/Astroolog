import documents from './documents.json'
export type LegalPage = 'privacy' | 'cookies' | 'terms'
export type LegalLocale = 'ru' | 'lv' | 'en'
export const legalDocuments = documents
export const legalPages: LegalPage[] = ['privacy', 'cookies', 'terms']
// Technical claims are maintained with the implementation, not editable marketing copy.
export const technicalCookieSections = [1, 2, 3, 4, 7]
export const legalKey = (page: LegalPage, index: number, field: string) => `legal.${page}.section${index + 1}.${field}`
// Only explicit list lines become lists; all content remains escaped Vue text.
export function legalBlocks(text: string) {
 const blocks: {type:'p'|'ul'; lines:string[]}[]=[]
 for (const paragraph of text.split(/\n\s*\n/)) {
  for(const line of paragraph.split('\n')) {
   const type=line.startsWith('- ')?'ul':'p'
   const last=blocks.at(-1)
   if(last?.type===type)last.lines.push(type==='ul'?line.slice(2):line)
   else blocks.push({type,lines:[type==='ul'?line.slice(2):line]})
  }
  // Preserve paragraph boundaries independently of consecutive block types.
  blocks.push({type:'p',lines:[]})
 }
 return blocks.filter(b=>b.lines.length)
}
