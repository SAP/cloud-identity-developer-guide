import { defineConfig } from 'vitepress';

export default defineConfig({
    title: 'AMS Client Libraries',
    description: 'Documentation for SAP Authorization Management Service (AMS) client libraries',
    themeConfig: {
        nav: [
            { text: 'Troubleshooting', link: '/Troubleshooting' },
            { text: 'Support', link: '/Support' }
        ],
        sidebar: [
            {
                text: 'Concepts',
                items: [
                    { text: 'Getting Started', link: '/Concepts/GettingStarted' },
                    { text: 'Authorization Checks', link: '/Concepts/AuthorizationChecks' },
                    { text: 'Testing', link: '/Concepts/Testing' },
                    { text: 'Technical Communication', link: '/Concepts/TechnicalCommunication' },
                    { text: 'Deploying DCL Policies', link: '/Concepts/DeployDCL' },
                    { text: 'ValueHelp', link: '/Concepts/ValueHelp' },
                    { text: 'Logging', link: '/Concepts/Logging' },
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
    }
});
