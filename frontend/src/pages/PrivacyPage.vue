<!-- TODO before production: confirm legal controller name, registration details, contact and retention periods. -->
<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
<template><article class="section container prose"><span class="eyebrow">{{ t('privacy.takingCareOfYourData') }}</span><h1>{{ t('siteLayout.privacy') }}</h1><h2>{{ t('privacy.whatInformationTheSiteReceives') }}</h2><p>{{ t('privacy.theFormsCollectYourNameAndEmail') }}</p><h2>{{ t('privacy.whyThisInformationIsNeeded') }}</h2><p>{{ t('privacy.toAnswerYourQuestionDiscussTheFormat') }}</p><h2>{{ t('privacy.howEnquiriesAreHandled') }}</h2><p>{{ t('privacy.yourEnquiryIsStoredInTheSites') }}</p><h2>{{ t('privacy.requestsAboutYourData') }}</h2><p>{{ t('privacy.toRequestACorrectionOrDeletionOf') }} <LocaleLink to="/contacts">{{ t('privacy.contactForm') }}</LocaleLink>{{ t('privacy.pleaseDontIncludeSensitiveInformationThatIsnt') }}</p><h2>{{ t('privacy.technicalInformation') }}</h2><p>{{ t('privacy.theCurrentSiteHasNoAdvertisingTrackers') }}</p><p>{{ t('privacy.languagePreference') }}</p></article></template>
