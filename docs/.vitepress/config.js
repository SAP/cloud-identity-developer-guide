import { defineConfig } from 'vitepress';

export default defineConfig({
    title: 'AMS Client Libraries',
    description: 'Documentation for SAP Authorization Management Service (AMS) client libraries',
    themeConfig: {
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
                    { text: 'Authorization Checks', link: '/concepts/AuthorizationChecks' },
                    { text: 'Testing', link: '/concepts/Testing' },
                    { text: 'Technical Communication', link: '/concepts/TechnicalCommunication' },
                    { text: 'Deploying DCL Policies', link: '/concepts/DeployDCL' },
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
            }
        ]
    },
    head: [
        ['link', { rel: 'icon', href: '/favicon.png', type: 'image/png' }]
    ]
});
