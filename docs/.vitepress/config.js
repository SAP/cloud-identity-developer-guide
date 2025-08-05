import languages from './languages';
import { defineConfig } from 'vitepress';

export default defineConfig({
    title: 'SCI Developer Guide',
    description: 'Application development with the SAP Cloud Identity Services client libraries',
    markdown: {
        languages,
        languageAlias: {
            cds: 'js'
        }
    },
    themeConfig: {
        search: {
            provider: 'local'
        },
        outline: {
            level: [2, 3], // Show both ## and ### headings in the TOC
        },
        nav: [
            { text: 'Samples', link: '/Samples' },
            { text: 'Troubleshooting', link: '/Troubleshooting' },
            { text: 'Support', link: '/Support' }
        ],
        sidebar: [
            {
                text: 'Authorization',
                items: [
                    { text: 'Getting Started', link: '/Authorization/GettingStarted' },
                    { text: 'Startup Check', link: '/Authorization/StartupCheck' },
                    { text: 'Authorization Checks', link: '/Authorization/AuthorizationChecks' },
                    { text: 'Testing', link: '/Authorization/Testing' },
                    { text: 'Technical Communication', link: '/Authorization/TechnicalCommunication' },
                    { text: 'Deploying DCL', link: '/Authorization/DeployDCL' },
                    { text: 'Value Help', link: '/Authorization/ValueHelp' },
                    { text: 'Logging', link: '/Authorization/Logging' },
                ]
            },
            {
                text: 'CAP Authorization',
                items: [
                    { text: 'Basics', link: '/CAP/Basics' },
                    { text: 'Instance-based Authorization', link: '/CAP/InstanceBasedAuthorization' },
                    { text: 'cds Plugin', link: '/CAP/cds-Plugin' }
                ]
            },
            {
                text: 'Library Modules',
                items: [
                    {
                        text: 'Java',
                        collapsed: true,
                        items: [
                            { text: 'jakarta-ams', link: '/Libraries/java/jakarta-ams/jakarta-ams' },
                            { text: 'spring-ams', link: '/Libraries/java/spring-ams/spring-ams' },
                            { text: 'cap-ams-support', link: '/Libraries/java/cap-ams-support/cap-ams-support' }
                        ]
                    },
                    {
                        text: 'Node.js',
                        collapsed: true,
                        items: [
                            { text: '@sap/ams', link: '/Libraries/nodejs/sap_ams/sap_ams' },
                            { text: '@sap/ams-dev', link: 'https://www.npmjs.com/package/@sap/ams-dev' }
                        ]
                    },
                    {
                        text: 'Go',
                        collapsed: true,
                        items: [
                            { text: 'cloud-identity-authorizations-golang-library', link: '/Libraries/go/go-ams' }
                        ]
                    }
                ]
            },
            {
                text: 'Resources',
                items: [
                    { text: 'Privacy', link: '/resources/Privacy' },
                    { text: 'Imprint', link: 'https://www.sap.com/about/legal/impressum.html', target: '_blank' },
                    { text: 'Terms of Use', link: 'https://www.sap.com/about/legal/terms-of-use.html', target: '_blank' },
                    { text: 'Trademarks', link: 'https://www.sap.com/about/legal/trademark.html', target: '_blank' }
                ]
            }
        ],
        footer: {
            message: '<a href="/resources/Privacy">Privacy</a> | <a href="https://www.sap.com/about/legal/impressum.html" target="_blank">Imprint</a> | <a href="https://www.sap.com/about/legal/terms-of-use.html" target="_blank">Terms of Use</a> | <a href="https://www.sap.com/about/legal/trademark.html" target="_blank">Trademarks</a>',
            copyright: '© 2025-present SAP SE or an SAP affiliate company and cloud-identity-authorizations-libraries contributors'
        },
        socialLinks: [
            { icon: 'github', link: 'https://github.com/SAP/cloud-identity-authorizations-libraries' },
        ]
    },
    head: [
        ['link', { rel: 'icon', href: '/favicon.png', type: 'image/png' }]
    ]
});
