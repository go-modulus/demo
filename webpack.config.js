const path = require("path")

module.exports = {
    entry: {
        "new-user": "./internal/user/page/js/new-user.js"
    },

    output: {
        filename: "[name].js",
        path: path.resolve(__dirname, "static/js")
    },

    mode: "production",
    devtool: "source-map",

    module: {
        rules: [
            {
                test: /\.js$/,
                exclude: [
                    /node_modules/
                ],
                use: [
                    { loader: "babel-loader" }
                ]
            }
        ]
    }
}