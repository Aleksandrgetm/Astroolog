export function validPhone(value: string): boolean {
 const phone=value.trim()
 if (!phone) return true
 return phone.length<=30 && /^\+?[0-9 ().-]+$/.test(phone) && /^\d{7,15}$/.test(phone.replace(/\D/g,''))
}
export function rigaToday(now=new Date()): string {
 const parts=new Intl.DateTimeFormat('en', {timeZone:'Europe/Riga',year:'numeric',month:'2-digit',day:'2-digit'}).formatToParts(now)
 const value=(type:string)=>parts.find(p=>p.type===type)!.value
 return `${value('year')}-${value('month')}-${value('day')}`
}
