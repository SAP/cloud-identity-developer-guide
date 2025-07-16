import { defineConfig } from 'vitepress';

export default defineConfig({
    title: 'AMS Client Libraries',
    description: 'Documentation for SAP Authorization Management Service (AMS) client libraries',
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
                text: 'Concepts',
                items: [
                    { text: 'Getting Started', link: '/concepts/GettingStarted' },
                    { text: 'Startup Check', link: '/concepts/StartupCheck' },
                    { text: 'Authorization Checks', link: '/concepts/AuthorizationChecks' },
                    { text: 'Testing', link: '/concepts/Testing' },
                    { text: 'Technical Communication', link: '/concepts/TechnicalCommunication' },
                    { text: 'Deploying DCL', link: '/concepts/DeployDCL' },
                    { text: 'ValueHelp', link: '/concepts/ValueHelp' },
                    { text: 'Logging', link: '/concepts/Logging' },
                ]
            },
            {
                text: 'CAP Integration',
                items: [
                    { text: 'Role Policies', link: '/CAP/RolePolicies' },
                    { text: 'Instance-based Authorization', link: '/CAP/InstanceBasedAuthorization' },
                    { text: 'DCL generation', link: '/CAP/DCLGeneration' }
                ]
            },
            {
                text: 'Java',
                items: [
                    { text: 'jakarta-ams', link: '/java/jakarta-ams/jakarta-ams' },
                    { text: 'spring-ams', link: '/java/spring-ams/spring-ams' },
                    { text: 'cap-ams-support', link: '/java/cap-ams-support/cap-ams-support' },
                    { text: 'cap-support (legacy)', link: '/java/cap-support/cap-support' }
                ]
            },
            {
                text: 'Node.js',
                items: [
                    { text: '@sap/ams', link: '/nodejs/sap_ams/sap_ams' },
                    { text: '@sap/ams-dev', link: '/nodejs/sap_ams-dev/sap_ams-dev' }
                ]
            },
            {
                text: 'Go',
                items: [
                    { text: 'cloud-identity-authorizations-golang-library', link: '/go/go-ams' }
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
