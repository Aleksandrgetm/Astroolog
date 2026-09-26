export const locales = ['ru','lv','en'] as const
export type SiteLocale = typeof locales[number]
export function localeFromPath(path:string):SiteLocale {return path==='/lv'||path.startsWith('/lv/')?'lv':path==='/en'||path.startsWith('/en/')?'en':'ru'}
export function basePath(path:string):string {return path.replace(/^\/(lv|en)(?=\/|$)/,'').replace(/\/$/,'') || '/'}
export function languagePath(path:string,locale:SiteLocale):string {const base=basePath(path);return locale==='ru'?base:`/${locale}${base==='/'?'/':base}`}
