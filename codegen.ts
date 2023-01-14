import type {CodegenConfig} from '@graphql-codegen/cli';

const config: CodegenConfig = {
    overwrite: true,
    schema: "http://localhost:8888/graphql",
    generates: {
        "integration/graphql.ts": {
            plugins: ["typescript"]
        },
        "./graphql.schema.json": {
            plugins: ["introspection"]
        }
    }
};

export default config;
