import { defineConfig } from 'vitepress';

export default defineConfig({
    title: 'AMS Client Libraries',
    description: 'Documentation for SAP Authorization Management Service (AMS) client libraries',
    themeConfig: {
        nav: [
            { text: 'Getting Started', link: '/GettingStarted' },
            { text: 'Authorization Checks', link: '/AuthorizationChecks' },
            { text: 'Testing', link: '/Testing' },
            { text: 'Technical Communication', link: '/TechincalCommunication' },
            { text: 'Deploying DCL Policies', link: '/DeployDCL' },
            { text: 'ValueHelp', link: '/ValueHelp' },
            { text: 'Support', link: '/Support' }
        ],
        sidebar: [
            {
                text: 'Concepts',
                items: [
                    { text: 'Getting Started', link: '/GettingStarted' },
                    { text: 'Authorization Checks', link: '/AuthorizationChecks' },
                    { text: 'Testing', link: '/Testing' },
                    { text: 'Technical Communication', link: '/TechincalCommunication' },
                    { text: 'Deploying DCL Policies', link: '/DeployDCL' },
                    { text: 'ValueHelp', link: '/ValueHelp' },
                    { text: 'Support', link: '/Support' }
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
    }
});
