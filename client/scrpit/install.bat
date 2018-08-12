# 快速安装node.js所需组件

# 进入项目根目录，安装vue等相关组件
yarn add vue vue-router element-ui whatwg-fetch query-string uglifyjs-webpack-plugin

# webpack.config.js
# 手动为file-loader增加ttf,woff,eot,woff2类型，否则将报错
# test: /\.(png|jpg|gif|svg|woff|ttf|eot|woff2)$/,
