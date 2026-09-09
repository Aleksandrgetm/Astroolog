import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'
import { createVueI18nAdapter } from 'vuetify/locale/adapters/vue-i18n'
import { useI18n } from 'vue-i18n'
import { i18n } from '../i18n'
import { VApp, VBtn, VIcon, VDialog, VForm, VTextField, VTextarea, VSelect, VCheckbox, VProgressCircular } from 'vuetify/components'
const components = { VApp, VBtn, VIcon, VDialog, VForm, VTextField, VTextarea, VSelect, VCheckbox, VProgressCircular }
import * as directives from 'vuetify/directives'
export default createVuetify({ locale: { adapter: createVueI18nAdapter({ i18n: i18n as unknown as Parameters<typeof createVueI18nAdapter>[0]['i18n'], useI18n }) }, components, directives, theme: { defaultTheme: 'warm', themes: { warm: { dark: false, colors: { primary: '#765239', secondary: '#b99573', background: '#faf7f2', surface: '#fffdfa', error: '#a33229' } } } }, defaults: { VBtn: { elevation: 0, rounded: 0 }, VTextField: { variant: 'outlined', color: 'primary' }, VTextarea: { variant: 'outlined', color: 'primary' }, VSelect: { variant: 'outlined', color: 'primary' } } })
