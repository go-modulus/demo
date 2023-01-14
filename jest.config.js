/** @type {import('ts-jest').JestConfigWithTsJest} */
module.exports = {
    preset: 'ts-jest',
    testEnvironment: 'node',
    setupFilesAfterEnv: [
        'jest-chain',
        'jest-extended/all',
    ],
    reporters: [
        'default',
        [
            '@jest-performance-reporter/core',
            {
                'errorAfterMs': 1000,
                'warnAfterMs': 500,
                'logLevel': 'warn',
                'maxItems': 5,
            },
        ]
    ],
};