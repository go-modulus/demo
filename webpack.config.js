const path = require("path")
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

module.exports = {
    plugins: [new MiniCssExtractPlugin({
        filename: '../css/[name].css',
    })],
    entry: {
        "new-user": "./internal/user/page/js/new-user.js",
        "blog/posts": "./internal/blog/page/js/posts.js",
        "blog": ["./internal/blog/page/css/blog.css"],
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
            },
            {
                test: /\.css$/i,
                use: [MiniCssExtractPlugin.loader, "css-loader"],
            },
        ]
    }
}